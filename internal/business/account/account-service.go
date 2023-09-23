package account

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/models"
)

type accountService struct {
	// depends on the interface not on the concrete imp
	accRepo core.AccountRepo
}

type AccountServiceConfig struct {
	AccRepo core.AccountRepo
}

func New(asc *AccountServiceConfig) core.AccountService {
	return &accountService{
		accRepo: asc.AccRepo,
	}
}

func (as *accountService) Create(ctx context.Context, reqDto *core.CreateAccReq) (*models.Account, error) {
	// convert request dto into a domain entity to pass it to the repo layer
	acc := new(models.Account)
	acc.OwnerName = reqDto.OwnerName
	acc.Currency = models.Currency(reqDto.Currency)
	acc.Balance = float64(0)
	acc.Removed = false

	// call the repo layer
	createdAcc, err := as.accRepo.Insert(ctx, acc)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}

	return createdAcc, nil
}

func (as *accountService) GetAll(ctx context.Context) ([]*models.Account, error) {
	accounts, err := as.accRepo.Get(ctx)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}
	return accounts, nil
}

func (as *accountService) GetByID(ctx context.Context, reqDto *core.GetAccByIdReq) (*models.Account, error) {
	accID := reqDto.ID
	acc, err := as.accRepo.GetByID(ctx, accID)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}
	return acc, nil
}

func (as *accountService) GetPage(ctx context.Context, reqDto *core.PaginateAccountsReq) ([]*models.Account, error) {
	// the service layer send to the repo layer the right calculated offset, its responsible for handling this logic
	limit := reqDto.Limit
	offset := (reqDto.Offset - 1) * limit
	accounts, err := as.accRepo.GetPage(ctx, limit, offset)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}
	return accounts, nil

}

func (as *accountService) UpdateByID(ctx context.Context, reqDto *core.UpdateAccountReq) (*models.Account, error) {
	updated, err := as.accRepo.Update(ctx, reqDto.ID, reqDto.Balance)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}
	return updated, nil
}

func (as *accountService) UpdateByOwnerName(ctx context.Context, reqDto *core.UpdateAccountByOwnerNameReq) (*models.Account, error) {
	updated, err := as.accRepo.UpdateByOwnerName(ctx, reqDto.OwnerName, reqDto.Balance)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return nil, err
	}
	return updated, nil
}

func (as *accountService) DeleteByID(ctx context.Context, reqDto *core.DeleteAccountReq) error {
	err := as.accRepo.Delete(ctx, reqDto.ID)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return err
	}
	return nil
}

func (as *accountService) DeleteByOwnerName(ctx context.Context, reqDto *core.DeleteAccountByOwnerNameReq) error {
	err := as.accRepo.DeleteByOwnerName(ctx, reqDto.OwnerName)
	if err != nil {
		log.Printf("[Account Service] | %v \n", err.Error())
		return err
	}
	return nil
}
