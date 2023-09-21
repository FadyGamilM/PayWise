package transfer

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/core"
	"paywise/internal/models"
)

type transferRepo struct {
	tx *sql.Tx
}

func New(tx *sql.Tx) core.TransferRepo {
	return &transferRepo{
		tx: tx,
	}
}

const (
	INSERT_TRANSFER_QUERY string = `
		INSERT INTO transfers 
		(to_account_id, from_account_id, amount)
		VALUES($1, $2, $3)
		RETURNING id
	`

	GET_BY_TRANSFER_ID_QUERY string = `
		SELECT id, to_account_id, from_account_id, amount 
		FROM transfers 
		WHERE id = $1
	`

	GET_ALL_TRANSFERS_FROM_ACCOUNT_QUERY_IN_PAGES = `
		SELECT id, to_account_id, from_account_id, amount 
		FROM transfers 
		WHERE  from_account_id = $1
		LIMIT $2
		OFFSET $3
	`

	GET_ALL_TRANSFERS_TO_ACCOUNT_QUERY_IN_PAGES = `
		SELECT id, to_account_id, from_account_id, amount 
		FROM transfers 
		WHERE  to_account_id = $1
		LIMIT $2
		OFFSET $3
	`

	GET_ALL_TRANSFERS_BETWEEN_TWO_ACCOUNTS_QUERY_IN_PAGES = `
		SELECT id, to_account_id, from_account_id, amount 
		FROM transfers 
		WHERE  to_account_id = $1 AND from_account_id = $2
		LIMIT $3
		OFFSET $4
	`
)

func (tr *transferRepo) Insert(ctx context.Context, transfer *models.Transfer) (int64, error) {
	var insertedTransferID int64
	if err := tr.tx.QueryRowContext(ctx, INSERT_TRANSFER_QUERY, transfer.ToAccountID, transfer.FromAccountID, transfer.Amount).Scan(&insertedTransferID); err != nil {
		log.Printf("error trying to scan the inserted transfer id => %v \n", err)
		return -1, err
	}

	return insertedTransferID, nil
}

func (tr *transferRepo) GetByID(ctx context.Context, transferID int64) (*models.Transfer, error) {
	transfer := new(models.Transfer)
	if err := tr.tx.QueryRowContext(ctx, GET_BY_TRANSFER_ID_QUERY, transferID).Scan(&transfer.ID, &transfer.ToAccountID, &transfer.FromAccountID, &transfer.Amount); err != nil {
		log.Printf("error trying to scan the retrieved transfer from database => %v \n", err)
	}

	return transfer, nil
}

func (tr *transferRepo) GetPageTransfersFromAcc(ctx context.Context, fromAccID int64, limit int16, offset int16) ([]*models.Transfer, error) {
	var transfers []*models.Transfer

	rows, err := tr.tx.QueryContext(ctx, GET_ALL_TRANSFERS_FROM_ACCOUNT_QUERY_IN_PAGES, fromAccID, limit, (offset-1)*limit)
	if err != nil {
		log.Printf("error trying to retrieve transfers page no.%v from account : %v to all accounts => %v \n", (offset - 1), fromAccID, err)
		return nil, err
	}
	for rows.Next() {
		transfer := new(models.Transfer)
		if err = rows.Scan(
			&transfer.ID,
			&transfer.ToAccountID,
			&transfer.FromAccountID,
			&transfer.Amount,
		); err != nil {
			log.Printf("error trying to scan the retrieved transfers rows from account : %v => %v \n", fromAccID, err)
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (tr *transferRepo) GetPageTransfersToAcc(ctx context.Context, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error) {
	var transfers []*models.Transfer

	rows, err := tr.tx.QueryContext(ctx, GET_ALL_TRANSFERS_TO_ACCOUNT_QUERY_IN_PAGES, toAccID, limit, (offset-1)*limit)
	if err != nil {
		log.Printf("error trying to retrieve transfers page no.%v to account : %v from all accounts => %v \n", (offset - 1), toAccID, err)
		return nil, err
	}
	for rows.Next() {
		transfer := new(models.Transfer)
		if err = rows.Scan(
			&transfer.ID,
			&transfer.ToAccountID,
			&transfer.FromAccountID,
			&transfer.Amount,
		); err != nil {
			log.Printf("error trying to scan the retrieved transfers rows to account : %v => %v \n", toAccID, err)
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (tr *transferRepo) GetPageTransfers(ctx context.Context, fromAccID int64, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error) {
	var transfers []*models.Transfer

	rows, err := tr.tx.QueryContext(ctx, GET_ALL_TRANSFERS_BETWEEN_TWO_ACCOUNTS_QUERY_IN_PAGES, toAccID, fromAccID, limit, (offset-1)*limit)
	if err != nil {
		log.Printf("error trying to retrieve transfers page no.%v from account : %v to account %v => %v \n", (offset - 1), fromAccID, toAccID, err)
		return nil, err
	}
	for rows.Next() {
		transfer := new(models.Transfer)
		if err = rows.Scan(
			&transfer.ID,
			&transfer.ToAccountID,
			&transfer.FromAccountID,
			&transfer.Amount,
		); err != nil {
			log.Printf("error trying to scan the retrieved transfers rows from account : %v to account : %v => %v \n", fromAccID, toAccID, err)
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
