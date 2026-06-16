package repository

import (
	"druna_server/pkg/model"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUserByUsername(username string) (model.User, error)
	GetUserByTelegramID(telegramID int64) (model.User, error)
	SearchUsers(prefix string) ([]model.FriendInfo, error)
}

type Token interface {
	RevokeToken(jti string, expiresAt time.Time) error
	IsTokenRevoked(jti string) (bool, error)
	Ping() error
}

type Event interface {
	CreateEvent(user model.Event) (int, error)
	UpdateEvent(userID int, event model.Event) error
	DeleteEvent(userID, eventID int) error
	HasOverlappingEvent(userID int, start, end time.Time, excludeID int) (bool, error)
	GetEventList(userID int) ([]model.Event, error)
	GetEventListFiltered(userID int, filter model.EventFilter) ([]model.Event, error)
	CountEvents(userID int, filter model.EventFilter) (int, error)
	GetEventsForUsers(userIDs []int, dateFrom, dateTo time.Time) (map[int][]model.Event, error)
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
	GetFriendshipStatus(userID, friendID int) (string, error)
	DeleteFriend(userID int, friendID int) error
}

type Group interface {
	CreateGroup(input model.Group) (int, error)
	ListGroups(userID int) ([]model.Group, error)
	GetGroupDetails(groupID, userID int) (model.GroupDetails, error)
	AddGroupMember(groupID, ownerID, memberID int) error
	DeleteGroup(groupID, ownerID int) error
	LeaveGroup(groupID, userID int) error
	ConfirmMemberTime(groupID, userID int, confirmedTime time.Time) error
	GetMemberUserIDs(groupID int) ([]int, error)
	IsGroupMember(groupID, userID int) (bool, error)
}

type Repository struct {
	Authorization
	Token
	Event
	Friendship
	Group
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Token:         NewTokenPostgres(db),
		Event:         NewEventPostgres(db),
		Friendship:    NewFriendshipPostgres(db),
		Group:         NewGroupPostgres(db),
	}
}
