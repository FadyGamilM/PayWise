package main

import (
	"context"
	"fmt"
	"log"
	"paywise/internal/database/postgres"
	accRepo "paywise/internal/repository/account"
	"time"
)

func main() {

	db, err := postgres.Setup()
	if err != nil {
		fmt.Errorf("error => ", err)
	}

	db.Ping()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// acc := &models.Account{
	// 	OwnerName: "menna",
	// 	Balance:   float64(200),
	// 	Currency:  models.EUR,
	// }

	// accID, err := accRepo.Insert(ctx, db, acc)
	// if err != nil {
	// 	fmt.Errorf("error ==> ", err)
	// }

	// log.Println("id : ", accID)

	account, err := accRepo.GetByID(ctx, db, 5)
	if err != nil {
		fmt.Errorf("error ==> ", err)
	}
	log.Println("id : ", account.ID)
	log.Println("ownername : ", account.OwnerName)
	log.Println("balance : ", account.Balance)

	// limit := 2
	// offset := 2
	// accounts, err := accRepo.GetPage(ctx, db, int16(limit), int16(offset-1)*int16(limit))
	// for _, acc := range accounts {
	// 	log.Println(acc.OwnerName)
	// }

	accRepo.Delete(ctx, db, 5)
	log.Println("deleted :D")

	account, err = accRepo.GetByID(ctx, db, 5)
	if err != nil {
		fmt.Errorf("error ==> ", err)
	}
	log.Println("id : ", account.ID)
	log.Println("ownername : ", account.OwnerName)
	log.Println("balance : ", account.Balance)
}
