package tests

import (
	"context"
	"database/sql"
	"net/url"
	"paywise/internal/repository/transactions"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	txStore = new(transactions.TxStore)
	tx      = new(sql.Tx)
)

func TestTransferMoneyTx(t *testing.T) {
	asserts := assert.New(t)

	dsn := url.URL{
		Scheme: "postgres",
		Host:   "localhost:2345",
		User:   url.UserPassword("paywise", "paywise"),
		Path:   "paywisedbtest",
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	var err error

	txStore.DB, err = sql.Open("pgx", dsn.String())
	if err != nil {
		t.Fatalf("error trying to open a postgres connection : %v", err.Error())
	}

	tx, err = txStore.DB.Begin()
	if err != nil {
		t.Fatalf("error trying to set the tx instance for test purpose : %v ", err.Error())
	}

	results := make(chan *transactions.TxTransferMoneyResult)
	errs := make(chan error)

	concurrent_Txs := 5
	amountToTransfer := float64(10)

	for i := 0; i < concurrent_Txs; i++ {
		go func() {
			txResult, err := txStore.TransferMoneyTX(context.Background(), &transactions.TxTransferMoneyArgs{
				FromAccountID: int64(1),
				ToAccountID:   int64(2),
				Amount:        amountToTransfer,
			})
			results <- txResult
			errs <- err
		}()
	}

	for i := 0; i < concurrent_Txs; i++ {
		txResult := <-results
		transfer := txResult.Transfer
		asserts.Equal(int64(1), transfer.FromAccountID)
		asserts.Equal(int64(2), transfer.ToAccountID)
		asserts.Equal(amountToTransfer, transfer.Amount)
		asserts.NotZero(transfer.ID)

		toEntry := txResult.ToEntry
		asserts.Equal(int64(2), toEntry.AccountID)
		asserts.Equal(amountToTransfer, toEntry.Amount)

		fromEntry := txResult.FromEntry
		asserts.Equal(int64(1), fromEntry.AccountID)
		asserts.Equal(-amountToTransfer, fromEntry.Amount)
		asserts.NotZero(fromEntry.ID)

	}

}

// type AccountRepoSuite struct {
// 	suite.Suite
// }

// var (
// 	txStore = new(transactions.TxStore)
// 	tx      = new(sql.Tx)
// 	fromAcc = new(models.Account)
// 	toAcc   = new(models.Account)
// )

// func TestAccountRepoSuite(t *testing.T) {
// 	suite.Run(t, &AccountRepoSuite{})
// }

// func (ars *AccountRepoSuite) SetupSuite() {
// 	ars.T().Log("setup the test suite environemnt ...")

// 	dsn := url.URL{
// 		Scheme: "postgres",
// 		Host:   "localhost:2345",
// 		User:   url.UserPassword("paywise", "paywise"),
// 		Path:   "paywisedbtest",
// 	}

// 	q := dsn.Query()
// 	q.Add("sslmode", "disable")

// 	dsn.RawQuery = q.Encode()

// 	dbInstance, err := sql.Open("pgx", dsn.String())
// 	if err != nil {
// 		ars.FailNowf("error trying to open a postgres connection", err.Error())
// 	}

// 	acc_1 := &models.Account{
// 		OwnerName: "fady",
// 		Balance:   float64(150),
// 		Currency:  models.EUR,
// 	}

// 	const INSERT_QUERY string = `
// 		INSERT INTO accounts (owner_name, balance, currency)
// 		VALUES ($1, $2, $3)
// 		RETURNING id, owner_name, balance, currency
// 	`

// 	if err := dbInstance.QueryRowContext(context.Background(), INSERT_QUERY, acc_1.OwnerName, acc_1.Balance, acc_1.Currency).Scan(&fromAcc.ID, &fromAcc.OwnerName, &fromAcc.Balance, &fromAcc.Currency); err != nil {
// 		ars.FailNow("error trying to isnert an account => %v \n", err)
// 	}

// 	acc_2 := &models.Account{
// 		OwnerName: "marwan",
// 		Balance:   float64(200),
// 		Currency:  models.EUR,
// 	}
// 	if err := dbInstance.QueryRowContext(context.Background(), INSERT_QUERY, acc_2.OwnerName, acc_2.Balance, acc_2.Currency).Scan(&toAcc.ID, &toAcc.OwnerName, &toAcc.Balance, &toAcc.Currency); err != nil {
// 		ars.FailNow("error trying to isnert an account => %v \n", err)
// 	}

// 	ars.T().Log("the from account => ", fromAcc.OwnerName, fromAcc.ID)
// 	ars.T().Log("the to account => ", toAcc.OwnerName, toAcc.ID)

// 	time.Sleep(10 * time.Second)

// 	txStore.DB = dbInstance

// 	tx, err = txStore.DB.Begin()
// 	if err != nil {
// 		ars.FailNowf("error trying to set the tx instance for test purpose", err.Error())
// 	}
// }

// func (ars *AccountRepoSuite) SetupTest() {
// 	ars.T().Log("setup before each unit test ....")
// 	accRepo := account.New(tx)

// 	acc_1 := &models.Account{
// 		OwnerName: "fady",
// 		Balance:   float64(150),
// 		Currency:  models.EUR,
// 	}

// 	createdFromAcc, err := accRepo.Insert(context.Background(), acc_1)
// 	if err != nil {
// 		ars.FailNowf("error trying to create the from account", err.Error())
// 	}

// 	acc_2 := &models.Account{
// 		OwnerName: "marwan",
// 		Balance:   float64(200),
// 		Currency:  models.EUR,
// 	}
// 	createdToAcc, err := accRepo.Insert(context.Background(), acc_2)
// 	if err != nil {
// 		ars.FailNowf("error trying to create the from account", err.Error())
// 	}

// 	fromAcc, err = accRepo.GetByID(context.Background(), createdFromAcc.ID)
// 	if err != nil {
// 		ars.FailNowf("error trying to fetch the from account", err.Error())
// 	}
// 	toAcc, err = accRepo.GetByID(context.Background(), createdToAcc.ID)
// 	if err != nil {
// 		ars.FailNowf("error trying to fetch the to account", err.Error())
// 	}

// 	ars.T().Log("the from account => ", fromAcc.OwnerName, fromAcc.ID)
// 	ars.T().Log("the to account => ", toAcc.OwnerName, toAcc.ID)
// }

// func (ars *AccountRepoSuite) TestCreateAccount() {
// 	ars.T().Log("running [test create account] ... ")

// 	results := make(chan *transactions.TxTransferMoneyResult)
// 	errs := make(chan error)

// 	concurrent_Txs := 5
// 	amountToTransfer := float64(10)

// 	ars.T().Log("the from account => ", fromAcc.OwnerName, fromAcc.ID)
// 	ars.T().Log("the to account => ", toAcc.OwnerName, toAcc.ID)

// 	for i := 0; i < concurrent_Txs; i++ {
// 		go func() {
// 			txResult, err := txStore.TransferMoneyTX(context.Background(), &transactions.TxTransferMoneyArgs{
// 				FromAccountID: fromAcc.ID,
// 				ToAccountID:   toAcc.ID,
// 				Amount:        amountToTransfer,
// 			})
// 			results <- txResult
// 			errs <- err
// 		}()
// 	}

// 	for i := 0; i < concurrent_Txs; i++ {
// 		err := <-errs
// 		ars.Nil(err)
// 		ars.NoError(err)

// 		result := <-results
// 		transfer := result.Transfer
// 		ars.Equal(fromAcc.ID, transfer.FromAccountID)
// 		ars.Equal(toAcc.ID, transfer.ToAccountID)
// 		ars.Equal(amountToTransfer, transfer.Amount)
// 		ars.NotZero(transfer.ID)

// 		toEntry := result.ToEntry
// 		ars.Equal(toAcc.ID, toEntry.AccountID)
// 		ars.Equal(amountToTransfer, toEntry.Amount)
// 		ars.NotZero(toEntry.ID)

// 		fromEntry := result.FromEntry
// 		ars.Equal(fromAcc.ID, fromEntry.AccountID)
// 		ars.Equal(-amountToTransfer, fromEntry.Amount)
// 		ars.NotZero(fromEntry.ID)
// 	}

// }
