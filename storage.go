package gobank

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	DeleteAccount(int) error
	GetAccountByNumber(int64) (*Account, error)
}

type PostgreesStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgreesStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreesStorage{
		db: db,
	}, nil
}

func (s *PostgreesStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgreesStorage) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		encrypted_password varchar(255),
		balance serial,
		created_at timestamp 
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgreesStorage) Seed() error {
	account, err := NewAccount("John", "Doe", "Hunter1234")
	if err != nil {
		return err
	}

	return s.CreateAccount(account)
}

func (s *PostgreesStorage) CreateAccount(account *Account) error {
	statement := `INSERT INTO account (first_name, last_name, number, encrypted_password, balance, created_at)
				  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(statement, account.FisrtName, account.LastName, account.Number, account.EncryptedPassword, account.Balance, account.CreatedAt)
	return err
}

func (s *PostgreesStorage) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgreesStorage) GetAccountByID(id int) (*Account, error) {

	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgreesStorage) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgreesStorage) GetAccountByNumber(accountNumber int64) (*Account, error) {

	rows, err := s.db.Query("SELECT * FROM account WHERE number = $1", accountNumber)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", accountNumber)
}

func (s *PostgreesStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FisrtName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}
