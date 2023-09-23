package transactions

import (
	"context"
	"log"
	"paywise/internal/core"
	transactionRepo "paywise/internal/repository/transactions"
)

type transactionService struct {
	txStore *transactionRepo.TxStore
}

type TransactionServiceConfig struct {
	TxStore *transactionRepo.TxStore
}

func (tsc *TransactionServiceConfig) New() core.TransactionService {
	return &transactionService{
		txStore: tsc.TxStore,
	}
}

func (ts *transactionService) TransferMoneyTransaction(ctx context.Context, reqDto *core.TxTransferMoneyReq) (*core.TxTransferMoneyRes, error) {
	result, err := ts.txStore.TransferMoneyTX(ctx, &transactionRepo.TxTransferMoneyArgs{
		ToAccountID:   reqDto.ToAccountID,
		FromAccountID: reqDto.FromAccountID,
		Amount:        reqDto.Amount,
	})
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return &core.TxTransferMoneyRes{
		Transfer:    result.Transfer,
		ToEntry:     result.ToEntry,
		FromEntry:   result.FromEntry,
		ToAccount:   result.ToAccount,
		FromAccount: result.FromAccount,
	}, nil
}
