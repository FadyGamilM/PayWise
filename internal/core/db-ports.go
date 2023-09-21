package core

import (
	"context"
	"paywise/internal/models"
)

type AccountRepo interface {
	Insert(ctx context.Context, acc *models.Account) (int64, error)
	Get(ctx context.Context) ([]*models.Account, error)
	GetByID(ctx context.Context, id int64) (*models.Account, error)
	GetPage(ctx context.Context, limit int16, offset int16) ([]*models.Account, error)
	Update(ctx context.Context, id int64, v float64) error
	UpdateByOwnerName(ctx context.Context, ownername string, v float64) error
	Delete(ctx context.Context, id int64) error
	DeleteByOwnerName(ctx context.Context, ownerName string) error
}

type EntryRepo interface {
	Insert(ctx context.Context, entry *models.Entry) (int64, error)
	Get(ctx context.Context, accID int64) ([]*models.Entry, error)
	GetbyID(ctx context.Context, accID int64, entryID int64) (*models.Entry, error)
	GetPage(ctx context.Context, accID int64, limit int16, offset int16) ([]*models.Entry, error)
}
