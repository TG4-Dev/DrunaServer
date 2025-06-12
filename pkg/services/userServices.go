package services

import (
	"BlobbyServer/pkg/models"
	"BlobbyServer/pkg/repositories"
	"BlobbyServer/pkg/utils"
	"errors"
)

var AuthService = authService{}

type authService struct{}

func (a *authService) Register(name, email, password string) (string, error) {
	exists, _ := repositories.UsersRepo.ExistsByEmail(email)
	if exists {
		return "", errors.New("user already exists")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}

	id, err := repositories.UsersRepo.Create(user)
	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(id)
}
