package repositories

import (
	"BlobbyServer/config"
	"BlobbyServer/pkg/models"
	"fmt"
	"time"
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
	err := config.DB.Exec(`
        INSERT INTO users (name, username, email, password_hash)
        VALUES (?, ?, ?, ?) RETURNING id
    `, user.Name, user.Username, user.Email, user.PasswordHash).Scan(&id).Error
	return id, err
}

func (r *usersRepo) SearchByEmail(email string) (models.User, error) {
	var results models.User
	err := config.DB.Raw(`
		SELECT id, password_hash 
		FROM users 
		WHERE email = ? 
		LIMIT 1
	`, email).Scan(&results).Error
	return results, err
}

func (r *usersRepo) ExistsByUsername(username string) (int, error) {
	var id int
	err := config.DB.Raw(`
		SELECT id FROM users WHERE username = ?
	`, username).Scan(&id).Error
	return id, err
}

func (r *usersRepo) CreateFriendRequest(userID int, friendID int) error {
	fmt.Println("Hello, World!")
	err := config.DB.Exec(`
		INSERT INTO friends (user_id, friend_id, status, request_at, confirmed_at)
		VALUES (?, ?, ?, ?, NULL)
	`, userID, friendID, "pending", time.Now()).Error
	return err
}
