package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"paywise/config"
	accountService "paywise/internal/business/account"
	authService "paywise/internal/business/auth"
	"paywise/internal/business/auth/paseto"
	entrytService "paywise/internal/business/entry"
	moneyTransactionService "paywise/internal/business/transactions"
	transferService "paywise/internal/business/transfer"
	userService "paywise/internal/business/user"
	"paywise/internal/database/postgres"
	accountRepo "paywise/internal/repository/account"
	entryRepo "paywise/internal/repository/entry"
	"paywise/internal/repository/transactions"
	transferRepo "paywise/internal/repository/transfer"
	userRepo "paywise/internal/repository/user"
	"paywise/internal/transport/rest"
	accountHandler "paywise/internal/transport/rest/handlers/account"
	authHandler "paywise/internal/transport/rest/handlers/auth"
	moneyTxHandler "paywise/internal/transport/rest/handlers/transactions"
	userHandler "paywise/internal/transport/rest/handlers/user"
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
	userRepo := userRepo.New(db)

	accServiceConfig := accountService.AccountServiceConfig{AccRepo: accRepo}
	accService := accountService.New(&accServiceConfig)

	configs, err := config.LoadPasetoTokenConfig("./config")
	if err != nil {
		fmt.Println("error trying to load config variables", err)
	}
	pasetoTokenAuth, err := paseto.New(configs.Paseto.SymmetricKey)
	if err != nil {
		log.Printf("error trying to create paseto token auth imp | %v \n", err)
	}
	authServiceConfig := authService.AuthServiceConfig{UserRepo: userRepo, TokenAuth: pasetoTokenAuth}
	authService := authService.New(&authServiceConfig)

	entryServiceConfig := entrytService.EntryServiceConfig{EntryRepo: entryRepo}
	_ = entrytService.New(&entryServiceConfig)

	transferServiceConfig := transferService.TransferServiceConfig{TransferRepo: transferRepo}
	_ = transferService.New(&transferServiceConfig)

	moneyTransactionServiceConfig := moneyTransactionService.TransactionServiceConfig{TxStore: txStore}
	moneyTransactionService := moneyTransactionService.New(&moneyTransactionServiceConfig)

	userServiceConfig := userService.UserServiceConfig{UserRepository: userRepo}
	userService := userService.New(&userServiceConfig)

	// create a router
	router := rest.CreateRouter()

	// inistantiate a handler to be up and waiting for mapping its methods to the routes when a request comes to one of them
	accountHandler.New(&accountHandler.AccountHandlerConfig{R: router, Service: accService})
	authHandler.New(&authHandler.AuthHandlerConfig{R: router, AuthService: authService})
	userHandler.New(&userHandler.UserHandlerConfig{R: router, UserService: userService})

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
