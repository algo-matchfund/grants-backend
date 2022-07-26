package watchdog

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/internal/config"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"golang.org/x/sync/semaphore"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	algorand_models "github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
)

type AlgorandWatchTask struct {
	ID        int64
	ProjectID string
	UserID    string
	Tx        algorand_models.Transaction
}
type AlgorandWatchdog struct {
	Watchdog

	client                  *indexer.Client
	algod                   *algod.Client
	watchSem                *semaphore.Weighted
	db                      *database.GrantsDatabase
	watchCh                 chan AlgorandWatchTask
	blockConfirmationWindow uint64
}

func NewAlgorandWatchdog(config *config.Config, db *database.GrantsDatabase) (*AlgorandWatchdog, error) {
	indexerClient, err := indexer.MakeClient(config.PaymentProviders.Algorand.Indexer.Address, config.PaymentProviders.Algorand.Indexer.Token)
	if err != nil {
		return nil, err
	}

	algodClient, err := algod.MakeClient(config.PaymentProviders.Algorand.Node.Address, config.PaymentProviders.Algorand.Node.Token)
	if err != nil {
		return nil, err
	}

	return &AlgorandWatchdog{
		client:                  indexerClient,
		algod:                   algodClient,
		db:                      db,
		watchCh:                 make(chan AlgorandWatchTask, config.PaymentProviders.Algorand.WatchQueue),
		watchSem:                semaphore.NewWeighted(int64(config.PaymentProviders.Algorand.WatchQueue)),
		blockConfirmationWindow: config.PaymentProviders.Algorand.BlockConfirmation,
	}, nil
}

func (w *AlgorandWatchdog) Init() bool {
	// up to 10 seconds for health check
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := w.client.HealthCheck().Do(ctx)
	if err != nil {
		log.Printf("AlgorandWatchdog.Init: Algorand client failed health check, error - %s\n", err)
		return false
	}

	if len(resp.Errors) != 0 {
		log.Println("AlgorandWatchdog.Init: Algorand indexer has errors:")
		for _, error := range resp.Errors {
			log.Printf("\t%s\n", error)
		}
		return false
	}

	// reindex known transactions, which are still marked as pending
	txs, err := w.db.GetPendingAlgorandTransactions()
	if err != nil {
		log.Println("AlgorandWatchdog.Init: failed to get pending Algorand transactions, reason: ", err)
		return false
	}

	log.Println("AlgorandWatchdog.Init: Reindexing...")
	wg := &sync.WaitGroup{}
	dbTx, err := w.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println("AlgorandWatchdog.Init: failed to start a DB transaction for indexing")
		return false
	}
	defer dbTx.Commit()

	for _, tx := range txs {
		wg.Add(1)
		go w.indexTransaction(wg, dbTx, tx)
	}

	wg.Wait()
	log.Println("AlgorandWatchdog.Init: Reindexing done")

	return true
}

func (w *AlgorandWatchdog) StartWatch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-w.watchCh:
				w.watchSem.Release(1)
				w.handleInput(task)
			}
		}
	}()
}

// Watch blocks execution until the watched value is put to the watching queue
func (w *AlgorandWatchdog) Watch(v interface{}) bool {
	task, ok := v.(AlgorandWatchTask)
	if !ok {
		log.Printf("AlgorandWatchdog.Watch: watch supports AlgorandWatchTask only, passed %v\n", v)
		return false
	}

	if err := w.watchSem.Acquire(context.Background(), 1); err != nil {
		log.Println("AlgorandWatchdog.Watch: watch failed, watch queue is full")
		return false
	}

	w.watchCh <- task
	return true
}

func (w *AlgorandWatchdog) WatchWithContext(ctx context.Context, v interface{}) bool {
	task, ok := v.(AlgorandWatchTask)
	if !ok {
		log.Printf("AlgorandWatchdog.WatchWithContext: watch supports AlgorandWatchTask only, passed %v\n", v)
		return false
	}

	if err := w.watchSem.Acquire(ctx, 1); err != nil {
		log.Println("AlgorandWatchdog.WatchWithContext: watch failed, watch queue is full")
		return false
	}

	w.watchCh <- task
	return true
}

func (w *AlgorandWatchdog) GetTransaction(txid string) (*algorand_models.Transaction, error) {
	resp, err := w.client.LookupTransaction(txid).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return &resp.Transaction, nil
}

func (w *AlgorandWatchdog) GetTransactionFromNode(txid string) (*algorand_models.Transaction, error) {
	resp, _, err := w.algod.PendingTransactionInformation(txid).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if resp.PoolError != "" {
		return nil, errors.New(resp.PoolError)
	}

	/* not a direct mapping */
	tx := &algorand_models.Transaction{
		Sender:      resp.Transaction.Txn.Sender.String(),
		Fee:         uint64(resp.Transaction.Txn.Fee),
		FirstValid:  uint64(resp.Transaction.Txn.FirstValid),
		LastValid:   uint64(resp.Transaction.Txn.LastValid),
		GenesisHash: resp.Transaction.Txn.GenesisHash[:],
		GenesisId:   resp.Transaction.Txn.GenesisID,
		Group:       resp.Transaction.Txn.Group[:],
		Id:          txid,
		Note:        resp.Transaction.Txn.Note,
		PaymentTransaction: algorand_models.TransactionPayment{
			Receiver: resp.Transaction.Txn.Receiver.String(),
			Amount:   uint64(resp.Transaction.Txn.Amount),
		},
		Signature: algorand_models.TransactionSignature{
			Sig: resp.Transaction.Sig[:],
		},
		Type:           string(resp.Transaction.Txn.Type),
		ConfirmedRound: resp.ConfirmedRound,
	}

	return tx, nil
}

func (w *AlgorandWatchdog) GetBlock(block uint64) (*algorand_models.Block, error) {
	resp, err := w.client.LookupBlock(block).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (w *AlgorandWatchdog) GetCurrentBlock() (uint64, error) {
	resp, err := w.client.HealthCheck().Do(context.Background())
	if err != nil {
		return 0, err
	}

	return resp.Round, nil
}

func (w *AlgorandWatchdog) getTxStatus(tx *algorand_models.Transaction, currentBlock uint64) string {
	if tx.ConfirmedRound != 0 && tx.ConfirmedRound+w.blockConfirmationWindow <= currentBlock {
		return models.AlgorandTransactionStatusConfirmed
	} else {
		return models.AlgorandTransactionStatusPending
	}
}

func (w *AlgorandWatchdog) updateIndex(dbtx *sql.Tx, idx int64, projectId, userId, status, confirmation_time string, tx *algorand_models.Transaction) {
	// update DB data
	err := w.db.UpdateAlgorandTransactionAtomic(dbtx, projectId, tx, status)
	if err != nil {
		dbtx.Rollback()
		log.Printf("AlgorandWatchdog#index%d: failed to update algorand tx %s status, reason: %s\n", idx, tx.Id, err)
		return
	}

	if status == models.AlgorandTransactionStatusConfirmed {
		// now that the transfer was confirmed, we can save it as project funding
		_, err = w.db.PostFundDetailedAtomic(dbtx, projectId, int64(tx.PaymentTransaction.Amount), userId, confirmation_time)
		if err != nil {
			dbtx.Rollback()
			log.Printf("AlgorandWatchdog#index%d: failed to update project %s funding based on tx %s, reason: %s\n", idx, projectId, tx.Id, err)
		}
	}
}

func (w *AlgorandWatchdog) indexTransaction(wg *sync.WaitGroup, dbtx *sql.Tx, tx *database.PendingAlgorandTransaction) {
	defer wg.Done()

	log.Printf("AlgorandWatchdog.index#%d: reindexing %s\n", tx.ID, tx.TxID)

	currentBlock, err := w.GetCurrentBlock()
	if err != nil {
		log.Printf("AlgorandWatchdog#index%d: failed to get the current block, reason: %s\n", tx.ID, err)
		return
	}

	// check that the tx is in the chain still (no chain reorg has happened) by looking it up once
	transaction, err := w.GetTransaction(tx.TxID)
	if err != nil {
		log.Printf("AlgorandWatchdog#index%d: transaction %s not found in indexer, trying node\n", tx.ID, tx.TxID)
		transaction, err = w.GetTransactionFromNode(tx.TxID)
		if err != nil {
			log.Printf("AlgorandWatchdog#index%d: transaction %s not found in node, marking as failed, reason: %s\n", tx.ID, tx.TxID, err)
			// tx was dropped, mark as failed
			w.db.MarkAlgorandTransactionAsFailedAtomic(dbtx, tx.TxID)
		}
	}
	status := w.getTxStatus(transaction, currentBlock)

	// Block with the tx is confirmed but status in DB is outdated
	if status == models.AlgorandTransactionStatusConfirmed && tx.Status != models.AlgorandTransactionStatusConfirmed {
		block, err := w.GetBlock(transaction.ConfirmedRound)
		if err != nil {
			log.Printf("AlgorandWatchdog#index%d: failed to get the transaction %s confirmation block time, reason: %s\n", tx.ID, tx.TxID, err)
			return
		}
		// update status and move to project contributions
		w.updateIndex(dbtx, tx.ID, tx.ProjectID, tx.UserID, models.AlgorandTransactionStatusConfirmed, time.Unix(int64(block.Timestamp), 0).UTC().Format(time.RFC3339), transaction)
	}

	// Tx not confirmed or tx is in a block but needs more confirmations
	if status == models.AlgorandTransactionStatusPending {
		// start the watchdog as usual
		go w.Watch(AlgorandWatchTask{
			ID:        0,
			ProjectID: tx.ProjectID,
			UserID:    tx.UserID,
			Tx:        *transaction,
		})
	}
}

func (w *AlgorandWatchdog) handleInput(task AlgorandWatchTask) {
	log.Printf("AlgorandWatchdog#%d: starting to track tx %s status\n", task.ID, task.Tx.Id)

	tx := task.Tx

	health, err := w.client.HealthCheck().Do(context.Background())
	if err != nil {
		log.Printf("AlgorandWatchdog#%d: indexer not available, exiting", task.ID)
		return
	}

	currentBlock := health.Round
	status := w.getTxStatus(&tx, currentBlock)
	newStatus := status

	// if tx has already confirmed or failed, no need to track
	if status == models.AlgorandTransactionStatusPending {
		for {
			tx, err := w.GetTransaction(task.Tx.Id)
			if err != nil {
				log.Printf("AlgorandWatchdog#%d: transaction %s not found in indexer, trying node\n", task.ID, task.Tx.Id)
				tx, err = w.GetTransactionFromNode(task.Tx.Id)
				if err != nil {
					log.Printf("AlgorandWatchdog#%d: transaction %s not found in node, reason: %s\n", task.ID, task.Tx.Id, err)
					return
				}
			}

			currentRound, err := w.GetCurrentBlock()
			if err != nil {
				log.Printf("AlgorandWatchdog#%d: cannot determine transaction %s status, trying again later, reason: %s\n", task.ID, task.Tx.Id, err)
				time.Sleep(time.Second * 5)
				continue
			}
			newStatus = w.getTxStatus(tx, currentRound)
			if newStatus != status {
				break
			}

			time.Sleep(time.Second * 5)
		}
	}

	log.Printf("AlgorandWatchdog#%d: transaction %s watch is done, saving results\n", task.ID, task.Tx.Id)

	writeTx, err := w.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Printf("AlgorandWatchdog#%d: failed to start a DB transaction, reason: %s\n", task.ID, err)
		return
	}

	// update DB data
	err = w.db.UpdateAlgorandTransactionAtomic(writeTx, task.ProjectID, &tx, newStatus)
	if err != nil {
		writeTx.Rollback()
		log.Printf("AlgorandWatchdog#%d: failed to update algorand tx %s status, reason: %s\n", task.ID, task.Tx.Id, err)
		return
	}

	if newStatus == models.AlgorandTransactionStatusConfirmed {
		block, err := w.GetBlock(tx.ConfirmedRound)
		if err != nil {
			log.Printf("AlgorandWatchdog#i%d: failed to get the transaction %s confirmation block time, reason: %s\n", task.ID, task.Tx.Id, err)
			return
		}
		// now that the transfer was confirmed, we can save it as project funding
		_, err = w.db.PostFundDetailedAtomic(writeTx, task.ProjectID, int64(tx.PaymentTransaction.Amount), task.UserID, time.Unix(int64(block.Timestamp), 0).UTC().Format(time.RFC3339))
		if err != nil {
			writeTx.Rollback()
			log.Printf("AlgorandWatchdog#%d: failed to update project %s funding based on tx %s, reason: %s\n", task.ID, task.ProjectID, task.Tx.Id, err)
		}
	}

	writeTx.Commit()
}
