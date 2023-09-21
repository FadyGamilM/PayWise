package entry

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/core"
	"paywise/internal/models"
)

const (
	INSERT_QUERY string = `
		INSERT INTO entries (account_id, amount)
		VALUES ($1, $2)
		RETURNING id, account_id, amount
	`

	GET_ALL_QUERY string = `
		SELECT id, account_id, amount
		FROM entries
		WHERE account_id = $1
	`

	GET_ENTRY_BY_ID string = `
		SELECT id , account_id, amount
		FROM entries
		WHERE account_id = $1, id = $2
	`
	GET_ENTRIES_IN_PAGES string = `
		SELECT id, account_id, amount 
		FROM entries 
		WHERE account_id = $1
		ORDER BY id
		LIMIT $2
		OFFSET $2
	`
)

type entryRepo struct {
	tx *sql.Tx
}

func New(tx *sql.Tx) core.EntryRepo {
	return &entryRepo{
		tx: tx,
	}
}

func (er *entryRepo) Insert(ctx context.Context, entry *models.Entry) (*models.Entry, error) {
	createdEntry := new(models.Entry)
	if err := er.tx.QueryRowContext(ctx, INSERT_QUERY, entry.AccountID, entry.Amount).Scan(&createdEntry.ID, &createdEntry.AccountID, &createdEntry.Amount); err != nil {
		log.Printf("error trying to isnert an entry => %v \n", err)
		return nil, err
	}

	return createdEntry, nil
}

func (er *entryRepo) Get(ctx context.Context, accID int64) ([]*models.Entry, error) {
	var entries []*models.Entry

	rows, err := er.tx.QueryContext(ctx, GET_ALL_QUERY, accID)
	if err != nil {
		log.Printf("error trying to retrieve all entries of account : %v => %v \n", accID, err)
		return nil, err
	}
	for rows.Next() {
		entry := new(models.Entry)
		if err = rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Amount,
		); err != nil {
			log.Printf("error trying to scan the retrieved entries rows for account : %v => %v \n", accID, err)
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (er *entryRepo) GetbyID(ctx context.Context, accID int64, entryID int64) (*models.Entry, error) {
	entry := new(models.Entry)
	if err := er.tx.QueryRowContext(ctx, GET_ENTRY_BY_ID, accID, entryID).Scan(&entry.ID, &entry.AccountID, &entry.Amount); err != nil {
		log.Printf("error trying to scan the retrieved entry with id : %v for account : %v  from database => %v \n", entryID, accID, err)
		return nil, err
	}

	return entry, nil
}

func (er *entryRepo) GetPage(ctx context.Context, accID int64, limit int16, offset int16) ([]*models.Entry, error) {
	var entries []*models.Entry

	rows, err := er.tx.QueryContext(ctx, GET_ENTRIES_IN_PAGES, accID, limit, (offset-1)*limit)
	if err != nil {
		log.Printf("error trying to retrieve entries page no.%v of account : %v => %v \n", (offset - 1), accID, err)
		return nil, err
	}
	for rows.Next() {
		entry := new(models.Entry)
		if err = rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Amount,
		); err != nil {
			log.Printf("error trying to scan the retrieved entries rows for account : %v => %v \n", accID, err)
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
