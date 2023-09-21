package account

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/core"
	"paywise/internal/models"
)

const (
	INSERT_QUERY string = `
		INSERT INTO accounts (owner_name, balance, currency) 
		VALUES ($1, $2, $3)
		RETURNING id, owner_name, balance, currency
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
		UPDATE accounts 
		SET removed = TRUE 
		WHERE id = $1
	`

	DELETE_BY_OWNER_NAME_QUERY string = `
		UPDATE accounts 
		SET removed = TRUE 
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

type accountRepo struct {
	tx *sql.Tx
}

func New(tx *sql.Tx) core.AccountRepo {
	return &accountRepo{tx: tx}
}

// TODO (1) => configure the options
// TODO (2) => build a database layer custom errors
func (ar *accountRepo) Insert(ctx context.Context, acc *models.Account) (*models.Account, error) {
	// the repo logic
	createdAcc := new(models.Account)
	if err := ar.tx.QueryRowContext(ctx, INSERT_QUERY, acc.OwnerName, acc.Balance, acc.Currency).Scan(&createdAcc.ID, &createdAcc.OwnerName, &createdAcc.Balance, &createdAcc.Currency); err != nil {
		log.Printf("error trying to isnert an account => %v \n", err)
		return nil, err
	}

	// return the result
	return createdAcc, nil
}

func (ar *accountRepo) Get(ctx context.Context) ([]*models.Account, error) {
	// the repo logic
	rows, err := ar.tx.QueryContext(ctx, GET_QUERY)
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

	// return the result
	return accounts, nil
}

func (ar *accountRepo) GetPage(ctx context.Context, limit int16, offset int16) ([]*models.Account, error) {
	// the repo logic
	rows, err := ar.tx.QueryContext(ctx, PAGINATE_QUERY, limit, offset)
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

	// return the result
	return accounts, nil
}

func (ar *accountRepo) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	// the repo logic
	account := new(models.Account)
	err := ar.tx.QueryRowContext(ctx, GET_BY_ID_QUERY, id).Scan(
		&account.ID,
		&account.OwnerName,
		&account.Balance,
		&account.Currency,
		&account.Removed,
	)
	if err != nil {
		log.Printf("error trying to scan the retrieved account from database => %v \n", err)
		return nil, err
	}

	// return the result
	return account, nil
}

func (ar *accountRepo) Update(ctx context.Context, id int64, v float64) error {
	_, err := ar.tx.ExecContext(ctx, UPDATE_BALANCE_BY_ID_QUERY, v, id)
	if err != nil {
		log.Printf("error trying to update the account => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func (ar *accountRepo) UpdateByOwnerName(ctx context.Context, ownername string, v float64) error {
	_, err := ar.tx.ExecContext(ctx, UPDATE_BALANCE_BY_OWNER_NAME_QUERY, v, ownername)
	if err != nil {
		log.Printf("error trying to update the account => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func (ar *accountRepo) Delete(ctx context.Context, id int64) error {

	_, err := ar.tx.ExecContext(ctx, DELETE_BY_ID_QUERY, id)
	if err != nil {
		log.Printf("error trying to delete the account => %v \n", err)
		return err
	}

	// return the result
	return nil
}

func (ar *accountRepo) DeleteByOwnerName(ctx context.Context, ownerName string) error {
	_, err := ar.tx.ExecContext(ctx, DELETE_BY_OWNER_NAME_QUERY, ownerName)
	if err != nil {
		log.Printf("error trying to delete the account => %v \n", err)
		return err
	}

	// return the result
	return nil
}
