package dtos

import "paywise/internal/models"

type TxTransferMoneyReq struct {
	ToAccountID   int64   `json:"to_account_id"`
	FromAccountID int64   `json:"from_account_id"`
	Amount        float64 `json:"amount"`
}

type TxTransferMoneyRes struct {
	Transfer    *models.Transfer `json:"transfer"`
	ToEntry     *models.Entry    `json:"to_entry"`
	FromEntry   *models.Entry    `json:"from_entry"`
	ToAccount   *models.Account  `json:"to_account"`
	FromAccount *models.Account  `json:"from_account"`
}
