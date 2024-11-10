package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(newAccount *Account) error
	GetAccountByID(id int) (*Account, error)
	GetAccounts() ([]*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error
	DepositToAccount(id int, amount int64) error
	WithdrawFromAccount(id int, amount int64) error
	TransferMoney(fromAccountId int, toAccountId int, amount int64) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
		balance float,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `insert into account
	(first_name, last_name, number, balance, created_at) 
	values ($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		a.FirstName,
		a.LastName,
		a.Number,
		a.Balance,
		a.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.GetAccountByID(id)
	if err != nil {
		return err
	}
	_, err = s.db.Query("delete from account where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT 
	id, 
	first_name, 
	last_name, 
	number, 
	balance, 
	created_at 
	FROM account`
	row, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	var accounts = []*Account{}
	for row.Next() {
		account, err := scanIntoAccount(row)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) DepositToAccount(id int, amount int64) error {
	_, err := s.GetAccountByID(id)
	if err != nil {
		return err
	}

	if amount < 0 {
		return fmt.Errorf("deposit amount invalid")
	}
	_, err = s.db.Query("UPDATE account SET balance = balance + $1 WHERE id = $2", amount, id)
	return err
}

func (s *PostgresStore) WithdrawFromAccount(id int, amount int64) error {
	account, err := s.GetAccountByID(id)
	if err != nil {
		return err
	}
	if amount < 0 {
		return fmt.Errorf("withdrawal amount invalid")
	}
	if account.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	_, err = s.db.Query("UPDATE account SET balance = balance - $1 WHERE id = $2", amount, id)
	return err
}

func (s *PostgresStore) TransferMoney(fromAccountId int, toAccountId int, amount int64) error {
	fromAccount, err := s.GetAccountByID(fromAccountId)
	if err != nil {
		return fmt.Errorf("sender account %d not found", fromAccountId)
	}

	_, err = s.GetAccountByID(toAccountId)
	if err != nil {
		return fmt.Errorf("receiver account %d not found", toAccountId)
	}

	if amount > fromAccount.Balance {
		return fmt.Errorf("insufficient funds in account %d", fromAccountId)
	}

	_, err = s.db.Query("UPDATE account SET balance = balance - $1 WHERE id = $2", amount, fromAccountId)
	if err != nil {
		return fmt.Errorf("failed to deduct amount from account %d", fromAccountId)
	}

	_, err = s.db.Query("UPDATE account SET balance = balance + $1 WHERE id = $2", amount, toAccountId)
	if err != nil {
		return fmt.Errorf("failed to transfer amount to account %d", toAccountId)
	}

	return nil
}

func scanIntoAccount(row *sql.Rows) (*Account, error) {
	account := &Account{}
	if err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	); err != nil {
		return nil, err
	}
	return account, nil
}
