package repository

import (
	"druna_server/pkg/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUserByUsername(username string) (model.User, error)
	GetUserByTelegramID(telegramID int64) (model.User, error)
}

type Event interface {
	CreateEvent(user model.Event) (int, error)
	DeleteEvent(userID, eventID int) error
	GetEventList(userID int) ([]model.Event, error)
}

type Friendship interface {
	CreateFriendRequest(userID, friendID int) error
	AcceptFriendRequest(userID int, friendID int) error
	RejectFriendRequest(userID int, friendID int) error
	ExistsByUsername(username string) (int, error)
	GetFriendList(userID int) ([]model.FriendInfo, error)
	GetIncomingFriendRequests(userID int) ([]model.FriendInfo, error)
	GetOutgoingFriendRequests(userID int) ([]model.FriendInfo, error)
	GetFriendRequestList(userID int) ([]model.FriendInfo, error)
	DeleteFriend(userID int, friendID int) error
}

type Group interface {
	CreateGroup(input model.Group) (int, error)
	ListGroups(userID int) ([]model.Group, error)
	GetGroupDetails(groupID, userID int) (model.GroupDetails, error)
	AddGroupMember(groupID, ownerID, memberID int) error
}

type Repository struct {
	Authorization
	Event
	Friendship
	Group
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Event:         NewEventPostgres(db),
		Friendship:    NewFriendshipPostgres(db),
		Group:         NewGroupPostgres(db),
	}
}
