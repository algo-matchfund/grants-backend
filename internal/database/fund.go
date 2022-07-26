package database

import (
	"database/sql"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"

	algorand_models "github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Fund struct {
	ProjectId string
	Amount    []uint64
}

func (db *GrantsDatabase) getFunds(matchingPoolId string) ([]*Fund, error) {
	var funds []*Fund
	query := db.builder.
		Select(`id`).
		From("projects")
	stmt, params := query.MustSql()
	rows, err := db.Query(stmt, params...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var projectId string
		err = rows.Scan(&projectId)

		if err != nil {
			log.Println(err)
			continue
		}

		var fund Fund
		fund.ProjectId = projectId

		// get funds for each project
		q := db.builder.
			Select("sum(amount)").
			From("contributions").
			Where("project_id = ?", projectId).
			Where("matching_round_id = ?", matchingPoolId).
			GroupBy("user_id")
		s, p := q.MustSql()
		r, err := db.Query(s, p...)

		if err != nil {
			log.Println(err)
			continue
		}

		for r.Next() {
			var amount uint64
			err = r.Scan(&amount)

			if err != nil {
				log.Println(err)
				continue
			}

			fund.Amount = append(fund.Amount, amount)
		}
		funds = append(funds, &fund)
	}
	return funds, nil
}

func (db *GrantsDatabase) preparePostFund(projectId string, amount int64, userId string) (string, []interface{}, error) {
	matchingRound, err := db.GetCurrentMatchingRound()
	if err != nil {
		return "", nil, err
	}

	stmt, params := db.builder.
		Insert("contributions").
		Columns("project_id", "matching_round_id", "amount", "user_id").
		Values(projectId, matchingRound.ID, amount, userId).
		Suffix("returning id").MustSql()

	if err != nil {
		return "", nil, err
	}

	return stmt, params, nil
}

func (db *GrantsDatabase) preparePostFundDetailed(projectId string, amount int64, userId string, contributionTime string) (string, []interface{}, error) {
	matchingRound, err := db.GetCurrentMatchingRound()
	if err != nil {
		return "", nil, err
	}

	stmt, params := db.builder.
		Insert("contributions").
		Columns("project_id", "matching_round_id", "amount", "user_id", "contribution_time").
		Values(projectId, matchingRound.ID, amount, userId, contributionTime).
		Suffix("returning id").MustSql()

	if err != nil {
		return "", nil, err
	}

	return stmt, params, nil
}

func (db *GrantsDatabase) PostFundAtomic(tx *sql.Tx, projectId string, amount int64, userId string) (*models.ProjectContributor, error) {
	contributor := models.ProjectContributor{
		Amount:        amount,
		ContributorID: &userId,
	}

	stmt, params, err := db.preparePostFund(projectId, amount, userId)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(stmt, params...).Scan(&contributor.ID)
	if err != nil {
		return nil, err
	}

	return &contributor, err
}

func (db *GrantsDatabase) PostFundDetailedAtomic(tx *sql.Tx, projectId string, amount int64, userId string, contributionTime string) (*models.ProjectContributor, error) {
	contributor := models.ProjectContributor{
		Amount:        amount,
		ContributorID: &userId,
	}

	stmt, params, err := db.preparePostFundDetailed(projectId, amount, userId, contributionTime)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(stmt, params...).Scan(&contributor.ID)
	if err != nil {
		return nil, err
	}

	return &contributor, err
}

func (db *GrantsDatabase) PostFund(projectId string, amount int64, userId string) (*models.ProjectContributor, error) {
	contributor := models.ProjectContributor{
		Amount:        amount,
		ContributorID: &userId,
	}

	stmt, params, err := db.preparePostFund(projectId, amount, userId)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(stmt, params...).Scan(&contributor.ID)
	if err != nil {
		return nil, err
	}

	return &contributor, err
}

func (db *GrantsDatabase) PostFundDetailed(projectId string, amount int64, userId string, contributionTime string) (*models.ProjectContributor, error) {
	contributor := models.ProjectContributor{
		Amount:        amount,
		ContributorID: &userId,
	}

	stmt, params, err := db.preparePostFundDetailed(projectId, amount, userId, contributionTime)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(stmt, params...).Scan(&contributor.ID)
	if err != nil {
		return nil, err
	}

	return &contributor, err
}

type PendingAlgorandTransaction struct {
	ID                int64
	TxID              string
	UserID            string
	ProjectID         string
	Status            string
	ConfirmationRound uint64
	Amount            uint64
}

func (db *GrantsDatabase) GetPendingAlgorandTransactions() ([]*PendingAlgorandTransaction, error) {
	stmt, params := db.builder.
		Select("id, txid, user_id, project_id, status, confirmation_round, amount").
		From("algorand_contributions").
		Where("confirmation_round is null or status='pending'").
		MustSql()

	rows, err := db.Query(stmt, params...)
	if err != nil {
		return nil, err
	}

	var transactions []*PendingAlgorandTransaction
	for rows.Next() {
		tx := &PendingAlgorandTransaction{}

		err := rows.Scan(&tx.ID, &tx.TxID, &tx.UserID, &tx.ProjectID, &tx.Status, &tx.ConfirmationRound, &tx.Amount)
		if err != nil {
			log.Println(err)
			continue
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (db *GrantsDatabase) GetAlgorandTransactions(projectId, userId string, pending bool) (*models.FundingTransactionsArray, error) {
	query := db.builder.
		Select("txid, sender_address, confirmation_round, fee, receiver_address, amount, signature, status").
		From("algorand_contributions").
		Where("user_id = ?", userId).
		Where("project_id = ?", projectId)

	if !pending {
		query.Where("confirmation_round is not null")
	}

	stmt, params := query.MustSql()

	rows, err := db.Query(stmt, params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := &models.FundingTransactionsArray{
		ProjectID: projectId,
		TxType:    "algorand",
	}

	var transactions []interface{}
	for rows.Next() {
		tx := &models.AlgorandTransaction{}

		var confirmationRound sql.NullInt64
		// not sure reading signature will work
		err = rows.Scan(&tx.ID, &tx.Sender, confirmationRound, &tx.Fee, &tx.Receiver, &tx.Amount, &tx.Signature, &tx.Status)
		if err != nil {
			log.Println(err)
			continue
		}

		if confirmationRound.Valid {
			tx.ConfirmedRound = confirmationRound.Int64
		}

		transactions = append(transactions, tx)
	}

	result.Transactions = transactions

	return result, nil
}

func (db *GrantsDatabase) SaveAlgorandTransaction(projectId, userId string, tx *algorand_models.Transaction, currentBlock uint64) (int64, error) {
	var status string
	if tx.ConfirmedRound != 0 && tx.ConfirmedRound+10 <= currentBlock {
		status = models.AlgorandTransactionStatusConfirmed
	} else if tx.ConfirmedRound != 0 {
		status = models.AlgorandTransactionStatusPending
	} else {
		status = models.AlgorandTransactionStatusError
	}

	stmt, params := db.builder.Insert("algorand_contributions").
		Columns("txid, sender_address, user_id, receiver_address, project_id, signature, fee, amount, confirmation_round, status").
		Values(tx.Id, tx.Sender, userId, tx.PaymentTransaction.Receiver, projectId, tx.Signature.Sig, tx.Fee, tx.PaymentTransaction.Amount, tx.ConfirmedRound, status).
		Suffix("returning id").
		MustSql()

	var id int64
	err := db.QueryRow(stmt, params...).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db *GrantsDatabase) prepareUpdateAlgorandTransaction(projectId string, tx *algorand_models.Transaction, status string) (stmt string, params []interface{}) {
	stmt, params = db.builder.Update("algorand_contributions").
		Set("fee", tx.Fee).
		Set("confirmation_round", tx.ConfirmedRound).
		Set("status", status).
		Where("project_id = ? and txid = ?", projectId, tx.Id).
		MustSql()

	return
}

func (db *GrantsDatabase) UpdateAlgorandTransactionAtomic(dbTx *sql.Tx, projectId string, tx *algorand_models.Transaction, status string) error {
	stmt, params := db.prepareUpdateAlgorandTransaction(projectId, tx, status)
	_, err := dbTx.Exec(stmt, params...)
	if err != nil {
		return err
	}

	return nil
}

func (db *GrantsDatabase) UpdateAlgorandTransaction(projectId string, tx *algorand_models.Transaction, status string) error {
	stmt, params := db.prepareUpdateAlgorandTransaction(projectId, tx, status)

	_, err := db.Exec(stmt, params...)
	if err != nil {
		return err
	}

	return nil
}

func (db *GrantsDatabase) MarkAlgorandTransactionAsFailedAtomic(tx *sql.Tx, txid string) error {
	stmt, params := db.builder.
		Update("algorand_contributions").
		Set("status", models.AlgorandTransactionStatusError).
		Where("txid = ?", txid).
		MustSql()

	_, err := tx.Exec(stmt, params...)
	return err
}

func (db *GrantsDatabase) MarkAlgorandTransactionAsFailed(txid string) error {
	stmt, params := db.builder.
		Update("algorand_contributions").
		Set("status", models.AlgorandTransactionStatusError).
		Where("txid = ?", txid).
		MustSql()

	_, err := db.Exec(stmt, params...)
	return err
}
