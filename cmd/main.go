package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"paywise/internal/database/postgres"
	"paywise/internal/transport/rest"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.Setup()
	if err != nil {
		fmt.Errorf("error => %v", err)
	}

	db.Ping()

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// create a router
	router := rest.CreateRouter()

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
