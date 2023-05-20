package main

import (
	"math/rand"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	AccountNumber int64  `json:"accountnumber"`
	Password      string `json:"password"`
}

type TransferRequest struct {
	ToAccount string `json:"toAccount"`
	Amount    string `json:"amount"`
}
type UserClaims struct {
	AccountNumber int64 `json:"accountnumber"`
	jwt.RegisteredClaims
}

type CreateAccountRequest struct {
	FisrtName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type Account struct {
	ID        int64     `json:"id"`
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
