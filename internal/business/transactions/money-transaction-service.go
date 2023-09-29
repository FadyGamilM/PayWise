package transactions

import (
	"context"
	"errors"
	"log"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	transactionRepo "paywise/internal/repository/transactions"
	"sync"
)

type transactionService struct {
	txStore *transactionRepo.TxStore
	accRepo core.AccountRepo
}

type TransactionServiceConfig struct {
	TxStore *transactionRepo.TxStore
	AccRepo core.AccountRepo
}

func New(tsc *TransactionServiceConfig) core.TransactionService {
	return &transactionService{
		txStore: tsc.TxStore,
		accRepo: tsc.AccRepo,
	}
}

func (ts *transactionService) TransferMoneyTransaction(ctx context.Context, reqDto *dtos.TxTransferMoneyReq) (*dtos.TxTransferMoneyRes, error) {
	// first we need to fetch the from-account balance to ensure that its more than or equalte the required amount to be transfered
	// define waitgroup for sync purpose, so we ensure that the transaction of the money-transfer doesn't start before the transaction of the account-repo finished
	var wg sync.WaitGroup
	isValidAmount := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		fromAccount, err := ts.accRepo.GetByID(ctx, reqDto.FromAccountID)
		if err != nil {
			log.Println("the error in fetching from-account to check its balance against the amount is => ", err)
			isValidAmount = false
		}
		if fromAccount.Balance >= reqDto.Amount {
			isValidAmount = true
		} else {
			isValidAmount = false
		}
	}()
	wg.Wait()

	if !isValidAmount {
		log.Println("amount is larger than the from-account-balance")
		return nil, errors.New("amount is larger than the from-account balance")
	}

	result, err := ts.txStore.TransferMoneyTX(ctx, &transactionRepo.TxTransferMoneyArgs{
		ToAccountID:   reqDto.ToAccountID,
		FromAccountID: reqDto.FromAccountID,
		Amount:        reqDto.Amount,
	})
	if err != nil {
		log.Printf("[Transaction Service] | %v \n", err.Error())
		return nil, err
	}

	return &dtos.TxTransferMoneyRes{
		Transfer:    result.Transfer,
		ToEntry:     result.ToEntry,
		FromEntry:   result.FromEntry,
		ToAccount:   result.ToAccount,
		FromAccount: result.FromAccount,
	}, nil
}
