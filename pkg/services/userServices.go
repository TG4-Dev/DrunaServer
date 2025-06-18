package services

import (
	"BlobbyServer/pkg/models"
	"BlobbyServer/pkg/repositories"
	"BlobbyServer/pkg/utils"
	"errors"
)

var AuthService = authService{}

type authService struct{}

func (a *authService) Register(name, username, email, password string) (string, error) {
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
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}

	id, err := repositories.UsersRepo.Create(user)
	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(id)
}

func (a *authService) Login(email, password string) (string, error) {
	results, err := repositories.UsersRepo.SearchByEmail(email)
	if err != nil {
		return "", err
	}

	if err := utils.ComparePassword(results.PasswordHash, password); err != nil {
		return "", err
	}

	return utils.GenerateJWT(results.ID)
}

func (a *authService) CheckJWTService(jwt string) error {
	err := utils.CheckJWT(jwt)
	return err
}
