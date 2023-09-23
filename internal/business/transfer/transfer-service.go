package transfer

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/models"
)

type transferService struct {
	// depend on the abstraction of the repo layer [**IMP**]
	transferRepo core.TransferRepo
}

type TransferServiceConfig struct {
	TransferRepo core.TransferRepo
}

func (tsc *TransferServiceConfig) New() core.TransferService {
	return &transferService{
		transferRepo: tsc.TransferRepo,
	}
}

func (ts *transferService) Create(ctx context.Context, reqDto *dtos.CreateTransferReq) (*models.Transfer, error) {
	transfer, err := ts.transferRepo.Insert(ctx, &models.Transfer{
		ToAccountID:   reqDto.ToAccountID,
		FromAccountID: reqDto.FromAccountID,
		Amount:        reqDto.Amount,
	})
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return transfer, nil
}

func (ts *transferService) GetByID(ctx context.Context, reqDto *dtos.GetTransferByIdReq) (*models.Transfer, error) {
	transfer, err := ts.transferRepo.GetByID(ctx, reqDto.TransferID)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return transfer, nil
}

func (ts *transferService) GetTransfersFromSpecificAccount(ctx context.Context, reqDto *dtos.GetTransfersFromAccountReq) ([]*models.Transfer, error) {
	transfers, err := ts.transferRepo.GetPageTransfersFromAcc(ctx, reqDto.FromAccountID, reqDto.Limit, (reqDto.Offset-1)*reqDto.Limit)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return transfers, nil
}

func (ts *transferService) GetTransfersToSpecificAccount(ctx context.Context, reqDto *dtos.GetTransfersToAccountReq) ([]*models.Transfer, error) {
	transfers, err := ts.transferRepo.GetPageTransfersToAcc(ctx, reqDto.ToAccountID, reqDto.Limit, (reqDto.Offset-1)*reqDto.Limit)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return transfers, nil
}

func (ts *transferService) GetPageTransfers(ctx context.Context, reqDto *dtos.GetTransfersBetweenTwoAccountsReq) ([]*models.Transfer, error) {
	transfers, err := ts.transferRepo.GetPageTransfers(ctx, reqDto.FromAccountID, reqDto.ToAccountID, reqDto.Limit, (reqDto.Offset-1)*reqDto.Limit)
	if err != nil {
		log.Printf("[SERVICE LAYER] | %v \n", err.Error())
		return nil, err
	}

	return transfers, nil
}
