package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB Database
)

type Database struct {
	Sqlx *sqlx.DB
}

// InitConn is a function used to initiate the connect with the database
func InitConn(dsn string) error {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Printf("Error connecting to db: %s", err)
		return err
	}

	DB.Sqlx = db

	return nil
}

// GetDBInstance gets the initalised instance of the database
func GetDBInstance() *Database {
	return &DB
}
