package db

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Email     string
	Password  string
	CreatedAt string
	ID        int
}

func CreateAccount(db *sqlx.DB, email, password string) (*Account, error) {
	// Hashear password — nunca guardamos texto plano
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	err = db.QueryRowx(
		`INSERT INTO accounts (email, password)
		 VALUES ($1, $2)
		 RETURNING id, email, password, created_at`,
		email, string(hash),
	).StructScan(account)

	return account, err
}

func GetAccountByEmail(db *sqlx.DB, email string) (*Account, error) {
	account := &Account{}
	err := db.QueryRowx(
		`SELECT id, email, password, created_at
		 FROM accounts WHERE email = $1`,
		email,
	).StructScan(account)
	return account, err
}

func ValidatePassword(account *Account, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	return err == nil
}
