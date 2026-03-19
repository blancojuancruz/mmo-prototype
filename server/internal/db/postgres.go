package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `db:"host"`
	User     string `db:"user"`
	Password string `db:"password"`
	DBName   string `db:"db_name"`
	Port     int    `db:"port"`
}

func NewPostgres(cfg Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("✅ PostgreSQL connected")
	return db
}
