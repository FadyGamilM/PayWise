package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"paywise/internal/models"
	accountRepository "paywise/internal/repository/account"
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
	txKey   = struct{}{}
	toAcc   = new(models.Account)
	fromAcc = new(models.Account)
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

	// i began this transaction because i need to create two accounts for testing purpose on a separate transaction to avoid mixing with the concurrent 5 transactions i am testing
	tx, err = txStore.DB.Begin()
	if err != nil {
		t.Fatalf("error trying to set the tx instance for test purpose : %v ", err.Error())
	}

	prepareAccounts := make(chan string)
	go func() {
		accRepo := accountRepository.New(tx)
		createdToAcc, err := accRepo.Insert(context.Background(), &models.Account{OwnerName: "mayar", Balance: float64(150), Currency: models.EUR})
		if err != nil {
			asserts.FailNowf("failed while creating a to account ", err.Error())
		}
		createdFromAcc, err := accRepo.Insert(context.Background(), &models.Account{OwnerName: "samy", Balance: float64(100), Currency: models.EUR})
		if err != nil {
			asserts.FailNowf("failed while creating a from account ", err.Error())
		}

		toAcc = createdToAcc
		fromAcc = createdFromAcc

		tx.Commit()

		prepareAccounts <- "done creating the two accounts"
	}()

	<-prepareAccounts

	log.Println("the to => ", toAcc.ID, " the from => ", fromAcc.ID)

	results := make(chan *transactions.TxTransferMoneyResult)
	errs := make(chan error)

	concurrent_Txs := 5
	amountToTransfer := float64(10)

	for i := 0; i < concurrent_Txs; i++ {
		ctx := context.WithValue(context.Background(), txKey, fmt.Sprintf("tx.(%v)", i))
		go func() {
			txResult, err := txStore.TransferMoneyTX(ctx, &transactions.TxTransferMoneyArgs{
				FromAccountID: fromAcc.ID,
				ToAccountID:   toAcc.ID,
				Amount:        amountToTransfer,
			})
			results <- txResult
			errs <- err
		}()
	}

	for i := 0; i < concurrent_Txs; i++ {
		txResult := <-results
		transfer := txResult.Transfer
		asserts.Equal(fromAcc.ID, transfer.FromAccountID) // 1
		asserts.Equal(toAcc.ID, transfer.ToAccountID)     // 2
		asserts.Equal(amountToTransfer, transfer.Amount)
		asserts.NotZero(transfer.ID)

		toEntry := txResult.ToEntry
		asserts.Equal(toAcc.ID, toEntry.AccountID)      // 2
		asserts.Equal(amountToTransfer, toEntry.Amount) // 1

		fromEntry := txResult.FromEntry
		asserts.Equal(fromAcc.ID, fromEntry.AccountID) // 1
		asserts.Equal(-amountToTransfer, fromEntry.Amount)
		asserts.NotZero(fromEntry.ID)

		toAccountAfterTransaction := txResult.ToAccount
		asserts.NotEmpty(toAccountAfterTransaction)
		asserts.Equal(toAccountAfterTransaction.ID, toAcc.ID)

		fromAccountAfterTransaction := txResult.FromAccount
		asserts.NotEmpty(fromAccountAfterTransaction)
		asserts.Equal(fromAccountAfterTransaction.ID, fromAcc.ID)

		// now we need to check the difference between the [toAcc (state before transaction)] and the toAccountAfterTransaction (state after transaction) and we expect to see that the difference in balance must be divisable by the amountToTransfer because each concurrent transaction will add 1*amountToTransfer so for 5 transaction, our toAccount will have balance = balance + 5*amount so its divisable by the amount
		balanceDifferenceOfToAccount := toAccountAfterTransaction.Balance - toAcc.Balance
		balanceDifferenceOfFromAccount := fromAcc.Balance - fromAccountAfterTransaction.Balance
		// must be same
		asserts.Equal(balanceDifferenceOfFromAccount, balanceDifferenceOfToAccount)

		t.Log("the balance difference of the to-account is : $", balanceDifferenceOfToAccount)
		t.Log("the balance difference of the from-account is : $", balanceDifferenceOfFromAccount)
		// asserts.Equal(balanceDifferenceOfFromAccount%amountToTransfer == 0)
		// asserts.Equal(balanceDifferenceOfFromAccount%amountToTransfer == 0)
	}

	dbInstance, err := sql.Open("pgx", dsn.String())
	if err != nil {
		t.Fatalf("error trying to open another postgres connection : %v", err.Error())
	}
	newTx, err := dbInstance.Begin()
	if err != nil {
		t.Fatalf("error trying to begin another transaction : %v", err.Error())
	}

	// after all conccurrent transactions, the final account balance must be decreased or increased by n * amount where n is the number of transactions
	// afterAllTransactions := make(chan string)
	toAccountAfterAllTransactions := new(models.Account)
	fromAccountAfterAllTransactions := new(models.Account)
	// go func() {

	accRepo := accountRepository.New(newTx)
	toAccountAfterAllTransactions, err = accRepo.GetByID(context.Background(), toAcc.ID)
	if err != nil {
		t.Log("error trying to fetch the to-account after all transactions to test its balance")
		t.FailNow()
	}

	fromAccountAfterAllTransactions, err = accRepo.GetByID(context.Background(), fromAcc.ID)
	if err != nil {
		t.Log("error trying to fetch the from-account after all transactions to test its balance")
		t.FailNow()
	}

	// afterAllTransactions <- "done fetching the accounts"
	// }()

	// <-afterAllTransactions

	asserts.Equal(toAcc.Balance+(float64(concurrent_Txs)*amountToTransfer), toAccountAfterAllTransactions.Balance)
	asserts.Equal(fromAcc.Balance-(float64(concurrent_Txs)*amountToTransfer), fromAccountAfterAllTransactions.Balance)

	newTx.Commit()

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
