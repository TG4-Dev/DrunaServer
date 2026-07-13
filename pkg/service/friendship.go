package service

import (
	"database/sql"
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
	"fmt"
)

type FriendshipService struct {
	repo         repository.Friendship
	authRepo     repository.Authorization
	notification *NotificationService
}

func NewFriendshipService(repo repository.Friendship, authRepo repository.Authorization, notification *NotificationService) *FriendshipService {
	return &FriendshipService{repo: repo, authRepo: authRepo, notification: notification}
}

func (s *FriendshipService) SendFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if friendID == userID {
		return errors.New("cannot send friend request to yourself")
	}

	status, err := s.repo.GetFriendshipStatus(userID, friendID)
	if err == nil {
		switch status {
		case "accepted":
			return errors.New("you are already friends")
		case "pending":
			return errors.New("friend request already pending")
		case "rejected":
			return errors.New("friend request was rejected; cannot send again")
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err := s.repo.CreateFriendRequest(userID, friendID); err != nil {
		return err
	}
	if s.notification != nil {
		if user, err := s.authRepo.GetUserByID(userID); err == nil {
			s.notification.EnqueueFriendRequest(friendID, user.Username)
		}
	}
	return nil
}

func (s *FriendshipService) AcceptFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	status, err := s.repo.GetFriendshipStatus(userID, friendID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("no friend request found")
	}
	if err != nil {
		return err
	}
	if status != "pending" {
		return fmt.Errorf("friend request is not pending (status: %s)", status)
	}

	return s.repo.AcceptFriendRequest(userID, friendID)
}

func (s *FriendshipService) RejectFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	status, err := s.repo.GetFriendshipStatus(userID, friendID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("no friend request found")
	}
	if err != nil {
		return err
	}
	if status != "pending" {
		return fmt.Errorf("friend request is not pending (status: %s)", status)
	}

	return s.repo.RejectFriendRequest(userID, friendID)
}

func (s *FriendshipService) FriendList(userID int) ([]model.FriendInfo, error) {
	return s.repo.GetFriendList(userID)
}

func (s *FriendshipService) FriendRequestList(userID int) ([]model.FriendInfo, error) {
	return s.repo.GetFriendRequestList(userID)
}

func (s *FriendshipService) IncomingFriendRequests(userID int) ([]model.FriendInfo, error) {
	return s.repo.GetIncomingFriendRequests(userID)
}

func (s *FriendshipService) OutgoingFriendRequests(userID int) ([]model.FriendInfo, error) {
	return s.repo.GetOutgoingFriendRequests(userID)
}

func (s *FriendshipService) DeleteFriend(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	return s.repo.DeleteFriend(userID, friendID)
}
