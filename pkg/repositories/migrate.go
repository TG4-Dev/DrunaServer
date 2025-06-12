package repositories

import (
	"BlobbyServer/config"
	"BlobbyServer/pkg/models"
)

func MigrateAll() {
	// Автомиграция - создание таблиц
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Event{},
		&models.GroupMember{},
		&models.Friend{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}
