package account

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/models"
)

const (
	INSERT_QUERY string = `
		INSERT INTO accounts (owner_name, balance, currency) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	GET_QUERY string = `
		SELECT id, owner_name, balance, currency, removed
		FROM accounts 
	`

	PAGINATE_QUERY string = `
		SELECT id, owner_name, balance, currency, removed
		FROM accounts 
		ORDER BY id 
		LIMIT $1
		OFFSET $2
	`

	GET_BY_OWNER_NAME_QUERY string = `
		SELECT id, owner_name, balance, currency, removed
		FROM accounts 
		WHERE owner_name = $1
	`

	GET_BY_ID_QUERY string = `
		SELECT id, owner_name, balance, currency, removed
		FROM accounts 
		WHERE id = $1
	`

	DELETE_BY_ID_QUERY string = `
		DELETE FROM accounts 
		WHERE id = $1
	`

	DELETE_BY_OWNER_NAME_QUERY string = `
		DELETE FROM accounts 
		WHERE owner_name = $1
	`

	UPDATE_BALANCE_BY_OWNER_NAME_QUERY string = `
		UPDATE accounts
		SET balance = $1
		WHERE owner_name = $2
	`

	UPDATE_BALANCE_BY_ID_QUERY string = `
		UPDATE accounts
		SET balance = $1
		WHERE id = $2
	`
)

// TODO (1) => configure the options
// TODO (2) => build a database layer custom errors
func Insert(ctx context.Context, db *sql.DB, acc *models.Account) (int64, error) {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return -1, err
	}

	// defer the rollback
	defer func() {
		_ = tx.Rollback()
	}()

	// the repo logic
	var insertedAccID int64
	if err = tx.QueryRowContext(ctx, INSERT_QUERY, acc.OwnerName, acc.Balance, acc.Currency).Scan(&insertedAccID); err != nil {
		log.Printf("error trying to isnert a user => %v \n", err)
		return -1, err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return -1, err
	}

	// return the result
	return insertedAccID, nil
}

func Get(ctx context.Context, db *sql.DB) ([]*models.Account, error) {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return nil, err
	}

	// defer the rollback
	defer func() {
		_ = tx.Rollback()
	}()

	// the repo logic
	rows, err := tx.QueryContext(ctx, GET_QUERY)
	if err != nil {
		log.Printf("error trying to fetch all accounts => %v \n", err)
		return nil, err
	}

	var accounts []*models.Account
	for rows.Next() {
		account := new(models.Account)
		err = rows.Scan(
			&account.ID,
			&account.OwnerName,
			&account.Balance,
			&account.Currency,
			&account.Removed,
		)
		if err != nil {
			log.Printf("error trying to scan the retrieved rows from database => %v \n", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return nil, err
	}

	// return the result
	return accounts, nil
}

func GetPage(ctx context.Context, db *sql.DB, limit int16, offset int16) ([]*models.Account, error) {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return nil, err
	}

	// defer the rollback
	defer func() {
		_ = tx.Rollback()
	}()

	// the repo logic
	rows, err := tx.QueryContext(ctx, PAGINATE_QUERY, limit, offset)
	if err != nil {
		log.Printf("error trying to fetch all accounts => %v \n", err)
		return nil, err
	}

	var accounts []*models.Account
	for rows.Next() {
		account := new(models.Account)
		err = rows.Scan(
			&account.ID,
			&account.OwnerName,
			&account.Balance,
			&account.Currency,
			&account.Removed,
		)
		if err != nil {
			log.Printf("error trying to scan the retrieved rows from database => %v \n", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return nil, err
	}

	// return the result
	return accounts, nil
}

func GetByID(ctx context.Context, db *sql.DB, id int64) (*models.Account, error) {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return nil, err
	}

	// defer rollback
	defer func() {
		_ = tx.Rollback()
	}()

	// the repo logic
	account := new(models.Account)
	err = tx.QueryRowContext(ctx, GET_BY_ID_QUERY, id).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Balance,
		&account.Currency,
		&account.Removed,
	)
	if err != nil {
		log.Printf("error trying to scan the retrieved row from database => %v \n", err)
		return nil, err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return nil, err
	}

	// return the result
	return account, nil
}

func Update(ctx context.Context, db *sql.DB, id int64, v float64) error {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return err
	}

	// defer rollback
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, UPDATE_BALANCE_BY_ID_QUERY, v, id)
	if err != nil {
		log.Printf("error trying to update the account => %v \n", err)
		return err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func UpdateByOwnerName(ctx context.Context, db *sql.DB, ownername string, v float64) error {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return err
	}

	// defer rollback
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, UPDATE_BALANCE_BY_OWNER_NAME_QUERY, v, ownername)
	if err != nil {
		log.Printf("error trying to update the account => %v \n", err)
		return err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func Delete(ctx context.Context, db *sql.DB, id int64) error {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return err
	}

	// defer rollback
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, DELETE_BY_ID_QUERY, id)
	if err != nil {
		log.Printf("error trying to delete the account => %v \n", err)
		return err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func DeleteByOwnerName(ctx context.Context, db *sql.DB, ownerName string) error {
	// start a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return err
	}

	// defer rollback
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, DELETE_BY_OWNER_NAME_QUERY, ownerName)
	if err != nil {
		log.Printf("error trying to delete the account => %v \n", err)
		return err
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return err
	}

	// return the result
	return nil
}
