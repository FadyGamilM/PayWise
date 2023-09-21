package models

import "time"

type Entry struct {
	id         int64     `json:"id"`
	account_id int64     `json:"account_id"`
	amount     float64   `json:"amount"`
	removed    bool      `json:"removed"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"created_at"`
}

func (e *Entry) GetID() int64 {
	return e.id
}

func (e *Entry) GetAccountID() int64 {
	return e.account_id
}

func (e *Entry) GetAmount() float64 {
	return e.amount
}

func (e *Entry) GetCreationDate() time.Time {
	return e.created_at
}

func (e *Entry) GetUpdateDate() time.Time {
	return e.updated_at
}

type EntryBuilder struct {
	id         int64     `json:"id"`
	account_id int64     `json:"account_id"`
	amount     float64   `json:"amount"`
	removed    bool      `json:"removed"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"created_at"`
}

func (eb *EntryBuilder) ID(v int64) {
	eb.id = v
}

func (eb *EntryBuilder) AccountID(v int64) {
	eb.account_id = v
}
func (eb *EntryBuilder) Amount(v float64) {
	eb.amount = v
}

func (eb *EntryBuilder) CreatedAt(v time.Time) {
	eb.created_at = v
}

func (eb *EntryBuilder) UpdatedAt(v time.Time) {
	eb.updated_at = v
}

func (eb *EntryBuilder) Build() *Entry {
	return &Entry{
		id:         eb.id,
		account_id: eb.account_id,
		amount:     eb.amount,
		removed:    eb.removed,
		created_at: eb.created_at,
		updated_at: eb.updated_at,
	}
}
