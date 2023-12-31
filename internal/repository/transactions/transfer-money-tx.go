package transactions

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/models"
	accountRepository "paywise/internal/repository/account"
	entryRepository "paywise/internal/repository/entry"
	transferRepository "paywise/internal/repository/transfer"
)

type TxTransferMoneyArgs struct {
	ToAccountID   int64   `json:"to_account_id"`
	FromAccountID int64   `json:"from_account_id"`
	Amount        float64 `json:"amount"`
}

type TxTransferMoneyResult struct {
	Transfer    *models.Transfer `json:"transfer"`
	ToEntry     *models.Entry    `json:"to_entry"`
	FromEntry   *models.Entry    `json:"from_entry"`
	ToAccount   *models.Account  `json:"to_account"`
	FromAccount *models.Account  `json:"from_account"`
}

var txKey = struct{}{}

func (txStore *TxStore) TransferMoneyTX(ctx context.Context, args *TxTransferMoneyArgs) (*TxTransferMoneyResult, error) {
	txResult := new(TxTransferMoneyResult)
	_ = txStore.execTransaction(ctx, func(tx *sql.Tx) error {
		// setup the repos to run the queries within this transaction instance
		accRepo := accountRepository.New(tx)
		entryRepo := entryRepository.New(tx)
		transferRepo := transferRepository.New(tx)
		var err error

		txNumber := ctx.Value(txKey)

		// define a result instance to update it through the transaction

		// create a transfer record for this money transaction operation
		log.Printf("[%v] | creating a transfer record \n", txNumber)
		transfer := new(models.Transfer)
		transfer.ToAccountID = args.ToAccountID
		transfer.FromAccountID = args.FromAccountID
		transfer.Amount = args.Amount
		txResult.Transfer, err = transferRepo.Insert(ctx, transfer)
		if err != nil {
			return err
		}

		// create a from-entry record which represents the money is withdrawn from the from-account
		log.Printf("[%v] | creating a from-entry record \n", txNumber)
		fromEntry := new(models.Entry)
		fromEntry.AccountID = args.FromAccountID
		fromEntry.Amount = -args.Amount
		txResult.FromEntry, err = entryRepo.Insert(ctx, fromEntry)
		if err != nil {
			return err
		}

		// create a to-entry record which represents the money
		log.Printf("[%v] | creating a to-entry record \n", txNumber)
		toEntry := new(models.Entry)
		toEntry.AccountID = args.ToAccountID
		toEntry.Amount = args.Amount
		txResult.ToEntry, err = entryRepo.Insert(ctx, toEntry)
		if err != nil {
			return err
		}

		// HINT : to avoid deadlock, we must ensure that our transaction acquired an exclusive lock on a correct order
		if args.ToAccountID < args.FromAccountID {
			log.Printf("[%v] | updating the to-account record \n", txNumber)
			txResult.ToAccount, err = accRepo.Update(ctx, args.ToAccountID, args.Amount)
			if err != nil {
				return err
			}
			log.Printf("[%v] | the to-account record after update is %v \n", txNumber, txResult.ToAccount.Balance)

			log.Printf("[%v] | updating the from-account record \n", txNumber)
			txResult.FromAccount, err = accRepo.Update(ctx, args.FromAccountID, -1*args.Amount)
			if err != nil {
				return err
			}
			log.Printf("[%v] | the from-account record after update is %v \n", txNumber, txResult.FromAccount.Balance)
		} else {
			log.Printf("[%v] | updating the from-account record \n", txNumber)
			txResult.FromAccount, err = accRepo.Update(ctx, args.FromAccountID, -1*args.Amount)
			if err != nil {
				return err
			}
			log.Printf("[%v] | the from-account record after update is %v \n", txNumber, txResult.FromAccount.Balance)

			log.Printf("[%v] | updating the to-account record \n", txNumber)
			txResult.ToAccount, err = accRepo.Update(ctx, args.ToAccountID, args.Amount)
			if err != nil {
				return err
			}
			log.Printf("[%v] | the to-account record after update is %v \n", txNumber, txResult.ToAccount.Balance)
		}

		return nil
	})
	return txResult, nil
}
