package main

import (
	"math/rand/v2"
	"time"
)

type DepositRequest struct {
	ToAccount int   `json:"toAccount"`
	Amount    int64 `json:"amount"`
}

type WithdrawRequest struct {
	FromAccount int   `json:"fromAccount"`
	Amount      int64 `json:"amount"`
}

type TransferRequest struct {
	FromAccount int   `json:"fromAccount"`
	ToAccount   int   `json:"toAccount"`
	Amount      int64 `json:"amount"`
}

type CreateAccountRequest struct {
	FistName string `json:"firstName"`
	LastName string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"Number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"time"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int64N(10000000),
		CreatedAt: time.Now().UTC(),
	}
}
