package transfer

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/database/postgres"
	"paywise/internal/models"
)

type transferRepo struct {
	pg *postgres.PG
}

func New(pg postgres.DBTX) core.TransferRepo {
	return &transferRepo{
		pg: &postgres.PG{DB: pg},
	}
}

const (
	INSERT_TRANSFER_QUERY string = `
		INSERT INTO transfers 
		(to_account, from_account, amount)
		VALUES($1, $2, $3)
		RETURNING id, to_account, from_account, amount
	`

	GET_BY_TRANSFER_ID_QUERY string = `
		SELECT id, to_account, from_account, amount 
		FROM transfers 
		WHERE id = $1
	`

	GET_ALL_TRANSFERS_FROM_ACCOUNT_QUERY_IN_PAGES = `
		SELECT id, to_account, from_account, amount 
		FROM transfers 
		WHERE  from_account = $1
		LIMIT $2
		OFFSET $3
	`

	GET_ALL_TRANSFERS_TO_ACCOUNT_QUERY_IN_PAGES = `
		SELECT id, to_account, from_account, amount 
		FROM transfers 
		WHERE  to_account = $1
		LIMIT $2
		OFFSET $3
	`

	GET_ALL_TRANSFERS_BETWEEN_TWO_ACCOUNTS_QUERY_IN_PAGES = `
		SELECT id, to_account, from_account, amount 
		FROM transfers 
		WHERE  to_account = $1 AND from_account = $2
		LIMIT $3
		OFFSET $4
	`
)

func (tr *transferRepo) Insert(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error) {
	createdTransfer := new(models.Transfer)
	if err := tr.pg.DB.QueryRowContext(ctx, INSERT_TRANSFER_QUERY, transfer.ToAccountID, transfer.FromAccountID, transfer.Amount).Scan(&createdTransfer.ID, &createdTransfer.ToAccountID, &createdTransfer.FromAccountID, &createdTransfer.Amount); err != nil {
		log.Printf("error trying to scan the inserted transfer id => %v \n", err)
		return nil, err
	}

	return createdTransfer, nil
}

func (tr *transferRepo) GetByID(ctx context.Context, transferID int64) (*models.Transfer, error) {
	transfer := new(models.Transfer)
	if err := tr.pg.DB.QueryRowContext(ctx, GET_BY_TRANSFER_ID_QUERY, transferID).Scan(&transfer.ID, &transfer.ToAccountID, &transfer.FromAccountID, &transfer.Amount); err != nil {
		log.Printf("error trying to scan the retrieved transfer from database => %v \n", err)
	}

	return transfer, nil
}

func (tr *transferRepo) GetPageTransfersFromAcc(ctx context.Context, fromAccID int64, limit int16, offset int16) ([]*models.Transfer, error) {
	var transfers []*models.Transfer

	rows, err := tr.pg.DB.QueryContext(ctx, GET_ALL_TRANSFERS_FROM_ACCOUNT_QUERY_IN_PAGES, fromAccID, limit, (offset-1)*limit)
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

	rows, err := tr.pg.DB.QueryContext(ctx, GET_ALL_TRANSFERS_TO_ACCOUNT_QUERY_IN_PAGES, toAccID, limit, (offset-1)*limit)
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

	rows, err := tr.pg.DB.QueryContext(ctx, GET_ALL_TRANSFERS_BETWEEN_TWO_ACCOUNTS_QUERY_IN_PAGES, toAccID, fromAccID, limit, (offset-1)*limit)
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
