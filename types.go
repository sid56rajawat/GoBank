package main

import "math/rand/v2"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"Number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.IntN(100000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int64N(10000000),
	}
}
