package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	accountService "paywise/internal/business/account"
	entrytService "paywise/internal/business/entry"
	moneyTransactionService "paywise/internal/business/transactions"
	transferService "paywise/internal/business/transfer"
	"paywise/internal/database/postgres"
	accountRepo "paywise/internal/repository/account"
	entryRepo "paywise/internal/repository/entry"
	"paywise/internal/repository/transactions"
	transferRepo "paywise/internal/repository/transfer"
	"paywise/internal/transport/rest"
	accountHandler "paywise/internal/transport/rest/handlers/account"
	moneyTxHandler "paywise/internal/transport/rest/handlers/transactions"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.Setup()
	if err != nil {
		log.Printf("error trying to connect to database : %v \n", err)
	}

	db.Ping()

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// setup transactions store and database repositories
	txStore := new(transactions.TxStore)
	txStore.DB = db

	accRepo := accountRepo.New(db)
	entryRepo := entryRepo.New(db)
	transferRepo := transferRepo.New(db)

	accServiceConfig := accountService.AccountServiceConfig{AccRepo: accRepo}
	accService := accountService.New(&accServiceConfig)

	entryServiceConfig := entrytService.EntryServiceConfig{EntryRepo: entryRepo}
	_ = entrytService.New(&entryServiceConfig)

	transferServiceConfig := transferService.TransferServiceConfig{TransferRepo: transferRepo}
	_ = transferService.New(&transferServiceConfig)

	moneyTransactionServiceConfig := moneyTransactionService.TransactionServiceConfig{TxStore: txStore}
	moneyTransactionService := moneyTransactionService.New(&moneyTransactionServiceConfig)

	// create a router
	router := rest.CreateRouter()

	// inistantiate a handler to be up and waiting for mapping its methods to the routes when a request comes to one of them
	accountHandler.New(&accountHandler.AccountHandlerConfig{R: router, Service: accService})

	moneyTxHandler.New(&moneyTxHandler.MoneyTxHandlerConfig{R: router, Service: moneyTransactionService})

	// create a server instance
	server := rest.CreateServer(router)

	// run the server up
	go rest.InitServer(server)

	// listen for shutdown or any interrupts
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for it
	<-quit
	// if we here, thats mean we will shut down the server gracefully
	rest.ShutdownGracefully(server)

}
