package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Cannot open DB connection:", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Cannot ping DB", err)
	}

	log.Println("Database connected successfully.")
}
