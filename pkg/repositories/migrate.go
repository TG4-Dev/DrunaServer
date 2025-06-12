package repositories

import (
	"BlobbyServer/pkg/models"

	"BlobbyServer/pkg/storage"
)

func MigrateAll() {
	db, err := storage.NewConnection()
	if err != nil {
		panic("failed to connect database")
	}

	// Автомиграция - создание таблиц
	err = db.AutoMigrate(
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
