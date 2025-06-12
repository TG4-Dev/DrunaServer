package repositories

import (
	"BlobbyServer/config"
	"BlobbyServer/pkg/models"
)

var UsersRepo = usersRepo{}

type usersRepo struct{}

func (r *usersRepo) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := config.DB.Raw(`SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`, email).Scan(&exists).Error
	return exists, err
}

func (r *usersRepo) Create(user models.User) (int, error) {
	var id int
	err := config.DB.Raw(`
        INSERT INTO users (name, email, password_hash)
        VALUES (?, ?, ?) RETURNING id
    `, user.Name, user.Email, user.PasswordHash).Scan(&id).Error
	return id, err
}
