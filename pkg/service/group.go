package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
)

type GroupService struct {
	repo     repository.Group
	friendRepo repository.Friendship
}

func NewGroupService(repo repository.Group, friendRepo repository.Friendship) *GroupService {
	return &GroupService{repo: repo, friendRepo: friendRepo}
}

func (s *GroupService) CreateGroup(input model.Group) (int, error) {
	if input.Name == "" {
		return 0, errors.New("group name is required")
	}
	return s.repo.CreateGroup(input)
}

func (s *GroupService) ListGroups(userID int) ([]model.Group, error) {
	return s.repo.ListGroups(userID)
}

func (s *GroupService) GetGroupDetails(groupID, userID int) (model.GroupDetails, error) {
	return s.repo.GetGroupDetails(groupID, userID)
}

func (s *GroupService) AddGroupMember(groupID, ownerID int, username string) error {
	memberID, err := s.friendRepo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	return s.repo.AddGroupMember(groupID, ownerID, memberID)
}
