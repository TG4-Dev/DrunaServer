package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"time"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(tokenTTL time.Duration, user model.User) (string, error)
	GenerateAccessRefreshToken(username, passwordHash string) (string, string, error)
	//GenerateRefreshToken(username, passwordHash string) (string, error)
	ParseToken(token string) (int, string, error)
	RenewToken(username string, userid int) (string, string, error)
}

type User interface {
}

type Event interface {
	CreateEvent(event model.Event) (int, error)
	DeleteEvent(userID, eventID int) error
	GetEventList(userID int) ([]model.Event, error)
}

type Friendship interface {
	SendFriendRequest(userID int, username string) error
	AcceptFriendRequest(userID int, username string) error
	RejectFriendRequest(userID int, username string) error
	FriendList(userID int) ([]model.FriendInfo, error)
	FriendRequestList(userID int) ([]model.FriendInfo, error)
	DeleteFriend(userID int, username string) error
}

type Group interface {
	CreateGroup(input model.Group) (int, error)
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
		Event:         NewEventService(repos.Event),
		Friendship:    NewFriendshipService(repos.Friendship),
		Group:         NewGroupService(repos.Group),
	}
}
