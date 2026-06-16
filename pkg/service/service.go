package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"time"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(tokenTTL time.Duration, user model.User) (string, error)
	GenerateAccessRefreshToken(username, password string) (string, string, error)
	ParseToken(token string) (int, string, error)
	RenewToken(username string, userid int) (string, string, error)
	TelegramLogin(telegramID int64, name, username string) (string, string, error)
	LoginWithTelegramInitData(initData string) (string, string, error)
}

type Event interface {
	CreateEvent(event model.Event) (int, error)
	DeleteEvent(userID, eventID int) error
	GetEventList(userID int) ([]model.Event, error)
	GetFreeTime(userID int, date time.Time) ([]model.TimeSlot, error)
}

type Friendship interface {
	SendFriendRequest(userID int, username string) error
	AcceptFriendRequest(userID int, username string) error
	RejectFriendRequest(userID int, username string) error
	FriendList(userID int) ([]model.FriendInfo, error)
	FriendRequestList(userID int) ([]model.FriendInfo, error)
	IncomingFriendRequests(userID int) ([]model.FriendInfo, error)
	OutgoingFriendRequests(userID int) ([]model.FriendInfo, error)
	DeleteFriend(userID int, username string) error
}

type Group interface {
	CreateGroup(input model.Group) (int, error)
	ListGroups(userID int) ([]model.Group, error)
	GetGroupDetails(groupID, userID int) (model.GroupDetails, error)
	AddGroupMember(groupID, ownerID int, username string) error
}

type Service struct {
	Authorization
	Event
	Friendship
	Group
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Event:         NewEventService(repos.Event),
		Friendship:    NewFriendshipService(repos.Friendship),
		Group:         NewGroupService(repos.Group, repos.Friendship),
	}
}
