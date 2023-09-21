package transactions

import (
	"context"
	"database/sql"
	"log"
)

type TxStore struct {
	DB *sql.DB
}

func (txStore *TxStore) execTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error {
	// start a transaction
	tx, err := txStore.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error trying to begin a transaction => %v \n", err)
		return err
	}

	// defer the rollback
	defer func() {
		_ = tx.Rollback()
	}()

	// execute a query withing the transaction
	if transactionErr := txFunc(tx); transactionErr != nil {
		log.Printf("error trying to execute the transaction => %v \n", transactionErr)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return commitErr
	}

	return nil
}
