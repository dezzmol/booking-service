package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func NewDB(host, username, password, dbname, port string) (*sqlx.DB, error) {
	if username == "" || password == "" || dbname == "" || port == "" {
		return nil, fmt.Errorf("[NewDB]: DB_USERNAME, DB_PASSWORD, DB_NAME или DB_PORT not initialized")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, dbname, port)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("[NewDB]: Error while connecting to db: %v", err)
	}

	return db, nil
}
