package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresStorage(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func InitStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database successfully connected!")
}
