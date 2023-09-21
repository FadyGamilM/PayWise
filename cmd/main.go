package main

import (
	"fmt"
	"log"
	"paywise/internal/database/postgres"
)

func main() {

	db, err := postgres.Setup()
	if err != nil {
		fmt.Errorf("error => ", err)
	}

	db.Ping()
	log.Println("done!")
}
