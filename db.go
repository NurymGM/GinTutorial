package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func openDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
		return err
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("DATABASE_URL is not set in the environment")
		return err
	}

	DB, err = sql.Open("postgres", connStr)
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