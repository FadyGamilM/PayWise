package account

import (
	"context"
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

func (as *accountService) Create(ctx context.Context, acc *models.Account) (*models.Account, error) {
	createdAcc, err := as.accRepo.Insert(ctx, acc)
	if err != nil {
		return 
	}
}

func (as *accountService) GetAll(ctx context.Context) ([]*models.Account, error) {}

func (as *accountService) GetByID(ctx context.Context, id int64) (*models.Account, error) {}

func (as *accountService) GetPage(ctx context.Context, limit int16, offset int16) ([]*models.Account, error) {
}

func (as *accountService) UpdateByID(ctx context.Context, id int64, v float64) (*models.Account, error) {
}

func (as *accountService) UpdateByOwnerName(ctx context.Context, ownername string, v float64) (*models.Account, error) {
}

func (as *accountService) DeleteByID(ctx context.Context, id int64) error {}

func (as *accountService) DeleteByOwnerName(ctx context.Context, ownerName string) error {}
