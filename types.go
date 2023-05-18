package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FisrtName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type Account struct {
	ID        int       `json:"id"`
	FisrtName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(fisrtName, lastName string) *Account {
	return &Account{
		FisrtName: fisrtName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		CreatedAt: time.Now().UTC(),
	}
}
