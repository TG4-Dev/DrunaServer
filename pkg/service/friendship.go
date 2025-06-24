package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
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

	err = s.repo.CreateFriendRequest(userID, friendID)
	if err != nil {
		return err
	}

	return err
}

func (s *FriendshipService) AcceptFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	err = s.repo.AcceptFriendRequest(userID, friendID)
	if err != nil {
		return err
	}

	return err
}

func (s *FriendshipService) RejectFriendRequest(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	err = s.repo.RejectFriendRequest(userID, friendID)
	if err != nil {
		return err
	}

	return err
}

func (s *FriendshipService) FriendList(userID int) ([]model.FriendInfo, error) {
	var friends []model.FriendInfo
	friends, err := s.repo.GetFriendList(userID)
	if err != nil {
		return []model.FriendInfo{}, err
	}

	return friends, err
}

func (s *FriendshipService) FriendRequestList(userID int) ([]model.FriendInfo, error) {
	var friends []model.FriendInfo
	friends, err := s.repo.GetFriendRequestList(userID)
	if err != nil {
		return []model.FriendInfo{}, err
	}

	return friends, err
}

func (s *FriendshipService) DeleteFriend(userID int, username string) error {
	friendID, err := s.repo.ExistsByUsername(username)
	if err != nil {
		return err
	}

	err = s.repo.DeleteFriend(userID, friendID)
	if err != nil {
		return err
	}

	return err
}
