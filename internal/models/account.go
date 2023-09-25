package models

import "time"

type Currency string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

type Account struct {
	ID        int64     `json:"id"`
	OwnerName string    `json:"owner_name"`
	Balance   float64   `json:"balance"`
	Currency  Currency  `json:"currency"`
	Removed   bool      `json:"removed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// func (a *Account) GetID() int64 {
// 	return a.id
// }

// func (a *Account) GetOwnerName() string {
// 	return a.owner_name
// }

// func (a *Account) GetBalance() float64 {
// 	return a.balance
// }

// func (a *Account) GetCurrency() Currency {
// 	return a.currency
// }

// func (a *Account) GetCreationDate() time.Time {
// 	return a.created_at
// }

// func (a *Account) GetUpdateDate() time.Time {
// 	return a.updated_at
// }

// type AccountBuilder struct {
// 	id         *int64     `json:"id"`
// 	owner_name *string    `json:"owner_name"`
// 	balance    *float64   `json:"balance"`
// 	currency   *Currency  `json:"currency"`
// 	removed    *bool      `json:"removed"`
// 	created_at *time.Time `json:"created_at"`
// 	updated_at *time.Time `json:"created_at"`
// }

// func (ab *AccountBuilder) ID(v *int64) {
// 	ab.id = v
// }

// func (ab *AccountBuilder) OwnerName(v *string) {
// 	ab.owner_name = v
// }
// func (ab *AccountBuilder) Balance(v *float64) {
// 	ab.balance = v
// }

// func (ab *AccountBuilder) Currency(v *Currency) {
// 	ab.currency = v
// }

// func (ab *AccountBuilder) IsRemoved(v *bool) {
// 	ab.removed = v
// }

// func (ab *AccountBuilder) CreatedAt(v *time.Time) {
// 	ab.created_at = v
// }

// func (ab *AccountBuilder) UpdatedAt(v *time.Time) {
// 	ab.updated_at = v
// }

// func (ab *AccountBuilder) Build() *Account {
// 	return &Account{
// 		id:         *ab.id,
// 		owner_name: *ab.owner_name,
// 		balance:    *ab.balance,
// 		currency:   *ab.currency,
// 		removed:    *ab.removed,
// 		created_at: *ab.created_at,
// 		updated_at: *ab.updated_at,
// 	}
// }
