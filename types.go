package gobank

import (
	"math/rand"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	AccountNumber int64  `json:"accountnumber"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	AccountNumber int64  `json:"accountnumber"`
	Token         string `json:"token"`
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
	Password  string `json:"password"`
}

type Account struct {
	ID                int64     `json:"id"`
	FisrtName         string    `json:"firstname"`
	LastName          string    `json:"lastname"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"_"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewAccount(fisrtName, lastName, password string) (*Account, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FisrtName:         fisrtName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(hashedPassword),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (a *Account) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
