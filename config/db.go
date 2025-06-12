package config

import (
	"log"
	"os"

	"gorm.io/gorm"

	"BlobbyServer/pkg/storage"
)

var DB *gorm.DB
var config storage.Config

func InitDB() {
	config.Host = os.Getenv("DB_HOST")
	config.Port = os.Getenv("DB_PORT")
	config.Password = os.Getenv("DB_PASSWORD")
	config.User = os.Getenv("DB_USER")
	config.DBName = os.Getenv("DB_NAME")
	config.SSLMode = os.Getenv("DB_SSLMODE")

	db, err := storage.NewConnection(config)
	DB = db

	if err != nil {
		log.Fatal("Cannot open DB connection:", err)
	}

	log.Println("Database connected successfully.")
}
