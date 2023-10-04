package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
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
	sessionRepo "paywise/internal/repository/session"
	"paywise/internal/repository/transactions"
	transferRepo "paywise/internal/repository/transfer"
	userRepo "paywise/internal/repository/user"
	"paywise/internal/transport/grpc/pb"
	_gRPC "paywise/internal/transport/grpc/server"
	"paywise/internal/transport/rest"
	accountHandler "paywise/internal/transport/rest/handlers/account"
	authHandler "paywise/internal/transport/rest/handlers/auth"
	moneyTxHandler "paywise/internal/transport/rest/handlers/transactions"
	userHandler "paywise/internal/transport/rest/handlers/user"
	"time"

	_ "github.com/lib/pq"
	gRPC "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//! 1. connect to postgres
	db, err := postgres.Setup()
	if err != nil {
		log.Printf("error trying to connect to database : %v \n", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("error trying to ping to the databse .. : %v", err)
		os.Exit(1)
	}

	log.Printf("Ponged successfully ..")
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//! 2. setup transactions store and database repositories
	txStore := new(transactions.TxStore)
	txStore.DB = db

	//! 3. setup repositories
	accRepo := accountRepo.New(db)
	entryRepo := entryRepo.New(db)
	transferRepo := transferRepo.New(db)
	userRepo := userRepo.New(db)
	sessionRepo := sessionRepo.New(db)

	//! 4. setup services
	pasetoConfigs, err := config.LoadPasetoTokenConfig("./config")
	if err != nil {
		fmt.Println("error trying to load config variables", err)
	}
	pasetoTokenAuth, err := paseto.New(pasetoConfigs.Paseto.SymmetricKey)
	if err != nil {
		log.Printf("error trying to create paseto token auth imp | %v \n", err)
	}
	accServiceConfig := accountService.AccountServiceConfig{AccRepo: accRepo, UserRepo: userRepo}
	accService := accountService.New(&accServiceConfig)

	authServiceConfig := authService.AuthServiceConfig{UserRepo: userRepo, SessionRepo: sessionRepo, TokenAuth: pasetoTokenAuth}
	authService := authService.New(&authServiceConfig)

	entryServiceConfig := entrytService.EntryServiceConfig{EntryRepo: entryRepo}
	_ = entrytService.New(&entryServiceConfig)

	transferServiceConfig := transferService.TransferServiceConfig{TransferRepo: transferRepo}
	_ = transferService.New(&transferServiceConfig)

	moneyTransactionServiceConfig := moneyTransactionService.TransactionServiceConfig{TxStore: txStore, AccRepo: accRepo}
	moneyTransactionService := moneyTransactionService.New(&moneyTransactionServiceConfig)

	userServiceConfig := userService.UserServiceConfig{UserRepository: userRepo}
	userService := userService.New(&userServiceConfig)

	//! 5.create a router
	router := rest.CreateRouter()

	//! 6. inistantiate a handler to be up and waiting for mapping its methods to the routes when a request comes to one of them
	accountHandler.New(&accountHandler.AccountHandlerConfig{R: router, Service: accService, TokenProvider: pasetoTokenAuth})
	authHandler.New(&authHandler.AuthHandlerConfig{R: router, AuthService: authService})
	userHandler.New(&userHandler.UserHandlerConfig{R: router, UserService: userService, TokenProvider: pasetoTokenAuth})
	moneyTxHandler.New(&moneyTxHandler.MoneyTxHandlerConfig{R: router, Service: moneyTransactionService, UserService: userService, TokenProvider: pasetoTokenAuth})

	//! ==> FOR HTTP REST BASED SERVER
	// // create a server instance
	// server := rest.CreateServer(router)

	// // run the server up
	// go rest.InitServer(server)

	grpcServer, err := _gRPC.NewServer(
		&_gRPC.GrpcServices{
			AuthService: authService,
		},
	)
	if err != nil {
		os.Exit(1)
	}
	log.Printf("the grpc address ==> %v", grpcServer.Address)
	gRPC_default_server := gRPC.NewServer()
	// register my implementation of grpc server which is (grpcServer)
	pb.RegisterPaywiseServer(gRPC_default_server, grpcServer)
	// register a reflection for our server so the clients can discover all the defined rpcs in our server
	reflection.Register(gRPC_default_server)
	// create a tcp listener for our grpc server
	listener, err := net.Listen("tcp", grpcServer.Address)
	if err != nil {
		log.Fatalf("error trying to create a tcp listener for a grpc server : %v", err)
	}
	log.Printf("start the grpc server on port : %v", listener.Addr().String())
	err = gRPC_default_server.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}

	// ! ==> FOR HTTP REST BASSED SERVER
	// // listen for shutdown or any interrupts
	// quit := make(chan os.Signal)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// // wait for it
	// <-quit
	// // if we here, thats mean we will shut down the server gracefully
	// rest.ShutdownGracefully(server)

}
