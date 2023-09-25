package models

import "time"

type Transfer struct {
	ID            int64     `json:"id"`
	ToAccountID   int64     `json:"to_account"`
	FromAccountID int64     `json:"from_account"`
	Amount        float64   `json:"amount"`
	Removed       bool      `json:"removed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// func (e *Transfer) GetID() int64 {
// 	return e.id
// }

// func (e *Transfer) GetToAccountID() int64 {
// 	return e.to_account_id
// }

// func (e *Transfer) GetFromAccountID() int64 {
// 	return e.from_account_id
// }

// func (e *Transfer) GetAmount() float64 {
// 	return e.amount
// }

// func (e *Transfer) GetCreationDate() time.Time {
// 	return e.created_at
// }

// func (e *Transfer) GetUpdateDate() time.Time {
// 	return e.updated_at
// }

// type TransferBuilder struct {
// 	id              int64     `json:"id"`
// 	to_account_id   int64     `json:"to_account_id"`
// 	from_account_id int64     `json:"from_account_id"`
// 	amount          float64   `json:"amount"`
// 	removed         bool      `json:"removed"`
// 	created_at      time.Time `json:"created_at"`
// 	updated_at      time.Time `json:"created_at"`
// }

// func (tb *TransferBuilder) ID(v int64) {
// 	tb.id = v
// }

// func (tb *TransferBuilder) ToAccountID(v int64) {
// 	tb.to_account_id = v
// }
// func (tb *TransferBuilder) FromAccountID(v int64) {
// 	tb.from_account_id = v
// }

// func (tb *TransferBuilder) Amount(v float64) {
// 	tb.amount = v
// }

// func (tb *TransferBuilder) CreatedAt(v time.Time) {
// 	tb.created_at = v
// }

// func (tb *TransferBuilder) UpdatedAt(v time.Time) {
// 	tb.updated_at = v
// }

// func (tb *TransferBuilder) Build() *Transfer {
// 	return &Transfer{
// 		id:              tb.id,
// 		to_account_id:   tb.to_account_id,
// 		from_account_id: tb.from_account_id,
// 		amount:          tb.amount,
// 		removed:         tb.removed,
// 		created_at:      tb.created_at,
// 		updated_at:      tb.updated_at,
// 	}
// }
