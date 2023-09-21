package transactions

import (
	"context"
	"database/sql"
	"log"
)

type TxStore struct {
	db *sql.DB
}

func (txStore *TxStore) ExecTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error {
	// start a transaction
	tx, err := txStore.db.BeginTx(ctx, nil)
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
	// now trying to rollback
	rollBackErr := tx.Rollback()
	if rollBackErr != nil {
		log.Printf("error trying to rollback the transaction => %v \n", rollBackErr)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("error trying to commit a transaction => %v \n", err)
		return commitErr
	}

	return nil
}

type TxTransferMoneyArgs struct {
}

type TxTransferMoneyResult struct {
}

func (txStore *TxStore) TransferMoneyTX(ctx context.Context, args *TxTransferMoneyArgs) (*TxTransferMoneyResult, error) {
	_ = txStore.ExecTransaction(ctx, func(tx *sql.Tx) error {
		panic("not implemented !")
	})
	return nil, nil
}
