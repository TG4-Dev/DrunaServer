package main

import (
	"context"
	druna "druna_server"
	"druna_server/pkg/handler"
	"druna_server/pkg/repository"
	"druna_server/pkg/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Druna API
// @version 1.0
// @description API server for Druna App

//@host localhost:8000
//@BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Warnf("No .env file loaded: %s", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("Failed to init DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	startTokenPurge(repos.Token)

	srv := new(druna.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err)
		}
	}()

	logrus.Println("DrunaServer started succesfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("DrunaServer shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutdown")
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db closing")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func startTokenPurge(tokenRepo repository.Token) {
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			count, err := tokenRepo.PurgeExpiredTokens()
			if err != nil {
				logrus.WithError(err).Error("failed to purge expired revoked tokens")
				continue
			}
			if count > 0 {
				logrus.WithField("count", count).Info("purged expired revoked tokens")
			}
		}
	}()
}
