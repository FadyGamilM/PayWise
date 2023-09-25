package core

import (
	"context"
	"paywise/internal/core/dtos"
	"paywise/internal/models"
)

type AuthService interface {
	Login(ctx context.Context, reqDto *dtos.LoginReq) (*dtos.LoginRes, error)
}

type UserRepo interface {
	Insert(ctx context.Context, user *models.User) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetAllAccounts(ctx context.Context, username string) ([]*models.Account, error)
}

type UserService interface {
	Create(ctx context.Context, reqDto *dtos.CreateUserDto) (*models.User, error)
	GetUserByUsername(ctx context.Context, reqDto *dtos.GetUserByUsernameDto) (*models.User, error)
	GetAllAccountsOfUserByUsername(ctx context.Context, reqDto *dtos.GetAllAccountsForUserDto) ([]*models.Account, error)
}

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
	Create(ctx context.Context, reqDto *dtos.CreateAccReq) (*models.Account, error)
	GetAll(ctx context.Context) ([]*models.Account, error)
	GetByID(ctx context.Context, reqDto *dtos.GetAccByIdReq) (*models.Account, error)
	GetPage(ctx context.Context, reqDto *dtos.PaginateAccountsReq) ([]*models.Account, error)
	UpdateByID(ctx context.Context, reqDto *dtos.UpdateAccountReq) (*models.Account, error)
	UpdateByOwnerName(ctx context.Context, reqDto *dtos.UpdateAccountByOwnerNameReq) (*models.Account, error)
	DeleteByID(ctx context.Context, reqDto *dtos.DeleteAccountReq) error
	DeleteByOwnerName(ctx context.Context, reqDto *dtos.DeleteAccountByOwnerNameReq) error
}

type EntryRepo interface {
	Insert(ctx context.Context, entry *models.Entry) (*models.Entry, error)
	Get(ctx context.Context, accID int64) ([]*models.Entry, error)
	GetbyID(ctx context.Context, accID int64, entryID int64) (*models.Entry, error)
	GetPage(ctx context.Context, accID int64, limit int16, offset int16) ([]*models.Entry, error)
}

type EntryService interface {
	Create(ctx context.Context, reqDto *dtos.CreateEntryReq) (*models.Entry, error)
	GetAll(ctx context.Context, reqDto *dtos.GetAllEntriesOfAccountReq) ([]*models.Entry, error)
	GetbyID(ctx context.Context, reqDto *dtos.GetEntryByIdReq) (*models.Entry, error)
	GetPage(ctx context.Context, reqDto *dtos.GetEntriesInPage) ([]*models.Entry, error)
}

type TransferRepo interface {
	Insert(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	GetByID(ctx context.Context, transferID int64) (*models.Transfer, error)
	GetPageTransfersFromAcc(ctx context.Context, fromAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetPageTransfersToAcc(ctx context.Context, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
	GetPageTransfers(ctx context.Context, fromAccID int64, toAccID int64, limit int16, offset int16) ([]*models.Transfer, error)
}

type TransferService interface {
	Create(ctx context.Context, reqDto *dtos.CreateTransferReq) (*models.Transfer, error)
	GetByID(ctx context.Context, reqDto *dtos.GetTransferByIdReq) (*models.Transfer, error)
	GetTransfersFromSpecificAccount(ctx context.Context, reqDto *dtos.GetTransfersFromAccountReq) ([]*models.Transfer, error)
	GetTransfersToSpecificAccount(ctx context.Context, reqDto *dtos.GetTransfersToAccountReq) ([]*models.Transfer, error)
	GetPageTransfers(ctx context.Context, reqDto *dtos.GetTransfersBetweenTwoAccountsReq) ([]*models.Transfer, error)
}

type TransactionService interface {
	TransferMoneyTransaction(ctx context.Context, reqDto *dtos.TxTransferMoneyReq) (*dtos.TxTransferMoneyRes, error)
}
