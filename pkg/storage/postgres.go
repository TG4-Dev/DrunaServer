package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection() (*gorm.DB, error) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		Password: "0851",
		User:     "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	dsn :=
		fmt.Sprintf(
			"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Password, config.User, config.DBName, config.SSLMode,
		)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, nil
}
