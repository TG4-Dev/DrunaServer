package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, passwordHash string) (string, error)
}

type User interface {
}

type Event interface {
}

type Friendship interface {
}

type Group interface {
}

type Service struct {
	Authorization
	User
	Event
	Friendship
	Group
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
