package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	DeleteAccount(int) (bool, error)
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
		balance serial,
		created_at timestamp 
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgreesStorage) CreateAccount(account *Account) error {
	statement := `INSERT INTO account (first_name, last_name, number, balance, created_at)
				  VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(statement, account.FisrtName, account.LastName, account.Number, account.Balance, account.CreatedAt)
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

	account := new(Account)

	for rows.Next() {

		err := rows.Scan(
			&account.ID,
			&account.FisrtName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

	}

	return account, nil
}

func (s *PostgreesStorage) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FisrtName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgreesStorage) DeleteAccount(id int) (bool, error) {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return true, nil
}
