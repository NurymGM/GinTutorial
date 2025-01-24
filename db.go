package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func openDB() error {
	var err error
	DB, err = sql.Open("postgres", "host=localhost port=5432 user=nurymalibekov dbname=library sslmode=disable")
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return err
	}
	return nil
}

func closeDB() error {
    if err := DB.Close(); err != nil {
        log.Printf("Failed to close database: %v", err)
        return err
    }
    return nil
}