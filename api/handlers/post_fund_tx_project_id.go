package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/service/watchdog"
	"github.com/go-openapi/runtime/middleware"

	algorand_models "github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type postNewProjectIDFundTxHandler struct {
	handler
	watchdogFactory *watchdog.WatchdogFactory
}

func (h *postNewProjectIDFundTxHandler) Handle(params operations.PostProjectIDFundTxParams, principal *models.Principal) middleware.Responder {
	log.Printf("POST /project/%s/fund/tx", params.ID)

	wd, err := h.watchdogFactory.GetWatchdog(params.Type)
	if err != nil {
		log.Printf("POST /project/%s/fund/tx: %s\n", params.ID, err)
		return operations.NewPostProjectIDFundTxInternalServerError()
	}

	switch params.Type {
	case "algorand":
		awd, ok := wd.(*watchdog.AlgorandWatchdog)
		if !ok {
			return operations.NewPostProjectIDFundTxInternalServerError()
		}
		return h.handleAlgorand(awd, params, principal)
	default:
		return operations.NewPostProjectIDFundTxBadRequest().WithPayload("Transaction type " + params.Type + " is not implemented")
	}
}

func (h *postNewProjectIDFundTxHandler) handleAlgorand(wd *watchdog.AlgorandWatchdog, params operations.PostProjectIDFundTxParams, principal *models.Principal) middleware.Responder {
	// check if transaction exists
	var tx *algorand_models.Transaction
	tx, err := wd.GetTransaction(params.Txid)
	if err != nil {
		log.Printf("POST /project/%s/fund/tx: transaction %s not found in indexer, trying node\n", params.ID, params.Txid)
		// maybe it's fresh, try to look up in the node
		tx, err = wd.GetTransactionFromNode(params.Txid)
		if err != nil {
			log.Printf("POST /project/%s/fund/tx: transaction %s not found, reason: %s\n", params.ID, params.Txid, err)
			return operations.NewPostProjectIDFundTxNotFound().WithPayload("Transaction " + params.Txid + " not found")
		}
	}

	projectWallet, err := h.db.GetAlgorandWalletByProjectId(params.ID)
	if err != nil {
		log.Printf("POST /project/%s/fund/tx: project %s not found or doesn't have a wallet, reason: %s\n", params.ID, params.ID, err)
		return operations.NewPostProjectIDFundTxNotFound().WithPayload("Project " + params.ID + " not found or doesn't have a wallet address")
	}

	if projectWallet != tx.PaymentTransaction.Receiver {
		return operations.NewPostProjectIDFundTxBadRequest().WithPayload("Receiver of transaction " + tx.Id + " is not project " + params.ID)
	}

	hash := sha256.New()
	hash.Write([]byte(params.ID))
	hash.Write([]byte(principal.ID))
	expected := hex.EncodeToString(hash.Sum(nil))

	// check if transaction was done by the user as a donation to the project
	if string(tx.Note) != expected {
		log.Printf("POST /project/%s/fund/tx: provided tx %s is not a donation to project %s (hash validation failed, expected \"%s\", got \"%s\")\n", params.ID, params.Txid, params.ID, expected, string(tx.Note))
		return operations.NewPostProjectIDFundTxNotFound().WithPayload("Transaction " + params.Txid + " is not a donation from you to the project")
	}

	currentBlock, err := wd.GetCurrentBlock()
	if err != nil {
		log.Printf("POST /project/%s/fund/tx: failed to get current Algorand block, reason: %s\n", params.ID, err)
		return operations.NewPostProjectIDFundTxInternalServerError()
	}

	id, err := h.db.SaveAlgorandTransaction(params.ID, principal.ID, tx, currentBlock)
	if err != nil {
		log.Printf("POST /project/%s/fund/tx: failed to save Algorand transaction, reason: %s\n", params.ID, err)
		return operations.NewPostProjectIDFundTxInternalServerError()
	}

	go wd.Watch(watchdog.AlgorandWatchTask{
		ID:        id,
		ProjectID: params.ID,
		UserID:    principal.ID,
		Tx:        *tx,
	})

	return operations.NewPostProjectIDFundTxOK()
}

// NewPostProjectIDFundTxHandler creates a handler for funding a project
func NewPostProjectIDFundTxHandler(db *database.GrantsDatabase, watchdogFactory *watchdog.WatchdogFactory) operations.PostProjectIDFundTxHandler {
	return &postNewProjectIDFundTxHandler{
		handler:         handler{db},
		watchdogFactory: watchdogFactory,
	}
}
