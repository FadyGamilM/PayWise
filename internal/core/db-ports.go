package core

import (
	"context"
	"paywise/internal/models"
)

type AccountRepo interface {
	Insert(ctx context.Context, acc *models.Account) (*models.Account, error)
	Get(ctx context.Context) ([]*models.Account, error)
	GetByID(ctx context.Context, id int64) (*models.Account, error)
	GetPage(ctx context.Context, limit int16, offset int16) ([]*models.Account, error)
	Update(ctx context.Context, id int64, v float64) (*models.Account, error)
	UpdateByOwnerName(ctx context.Context, ownername string, v float64) (*models.Account, error)
	Delete(ctx context.Context, id int64) error
	DeleteByOwnerName(ctx context.Context, ownerName string) error
}

type AccountService interface {
	Create(ctx context.Context, reqDto *CreateAccReq) (*models.Account, error)
	GetAll(ctx context.Context) ([]*models.Account, error)
	GetByID(ctx context.Context, reqDto *GetAccByIdReq) (*models.Account, error)
	GetPage(ctx context.Context, reqDto *PaginateAccountsReq) ([]*models.Account, error)
	UpdateByID(ctx context.Context, reqDto *UpdateAccountReq) (*models.Account, error)
	UpdateByOwnerName(ctx context.Context, reqDto *UpdateAccountByOwnerNameReq) (*models.Account, error)
	DeleteByID(ctx context.Context, reqDto *DeleteAccountReq) error
	DeleteByOwnerName(ctx context.Context, reqDto *DeleteAccountByOwnerNameReq) error
}

type EntryRepo interface {
	Insert(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	Get(ctx context.Context, accID int64) ([]*models.Entry, error)
	GetbyID(ctx context.Context, accID int64, entryID int64) (*models.Entry, error)
	GetPage(ctx context.Context, accID int64, limit int16, offset int16) ([]*models.Entry, error)
}

type EntryService interface {
	Create(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	GetAll(ctx context.Context, accID int64) ([]*models.Entry, error)
	GetbyID(ctx context.Context, accID int64, entryID int64) (*models.Entry, error)
	GetPage(ctx context.Context, accID int64, limit int16, offset int16) ([]*models.Entry, error)
}

type TransferRepo interface {
	Insert(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	GetByID(ctx context.Context, transferID int64) (*models.Transfer, error)
	GetPageTransfersFromAcc(ctx context.Context, fromAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetPageTransfersToAcc(ctx context.Context, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetPageTransfers(ctx context.Context, fromAccID int64, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
}

type TransferService interface {
	Create(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	GetByID(ctx context.Context, transferID int64) (*models.Transfer, error)
	GetTransfersFromSpecificAccount(ctx context.Context, fromAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetTransfersToSpecificAccount(ctx context.Context, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetPageTransfers(ctx context.Context, fromAccID int64, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
}

type TransactionService interface {
	TransferMoneyTransaction(ctx context.Context, reqDto *TxTransferMoneyReq) (*TxTransferMoneyRes, error)
}
