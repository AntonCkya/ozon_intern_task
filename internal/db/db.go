package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func InitDB(cfg DBConfig) (*sql.DB, error) {
	uri := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
