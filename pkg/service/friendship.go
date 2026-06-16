package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
)

type FriendshipService struct {
	repo repository.Friendship
}

func NewFriendshipService(repo repository.Friendship) *FriendshipService {
	return &FriendshipService{repo: repo}
}

func (s *FriendshipService) SendFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if friendID == userID {
		return errors.New("cannot send friend request to yourself")
	}

	return s.repo.CreateFriendRequest(userID, friendID)
}

func (s *FriendshipService) AcceptFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	return s.repo.AcceptFriendRequest(userID, friendID)
}

func (s *FriendshipService) RejectFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
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
