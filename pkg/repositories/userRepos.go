package repositories

import (
	"BlobbyServer/config"
	"BlobbyServer/pkg/models"
)

var UsersRepo = usersRepo{}

type usersRepo struct{}

func (r *usersRepo) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *usersRepo) Create(user models.User) (int, error) {
	var id int
	err := config.DB.QueryRow(`
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3) RETURNING id
	`, user.Name, user.Email, user.PasswordHash).Scan(&id)
	return id, err
}
