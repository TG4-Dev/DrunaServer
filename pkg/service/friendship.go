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

func (s *FriendshipService) FriendRequest(userID int, username string) error {
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

func (s *FriendshipService) FriendList(userID int) ([]model.FriendInfo, error) {
	var friends []model.FriendInfo
	friends, err := s.repo.GetFriendList(userID)
	if err != nil {
		return []model.FriendInfo{}, err
	}

	return friends, err
}
