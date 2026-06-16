package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
	"time"
)

type GroupService struct {
	repo      repository.Group
	friendRepo repository.Friendship
	eventRepo repository.Event
}

func NewGroupService(repo repository.Group, friendRepo repository.Friendship, eventRepo repository.Event) *GroupService {
	return &GroupService{repo: repo, friendRepo: friendRepo, eventRepo: eventRepo}
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
	isMember, err := s.repo.IsGroupMember(groupID, memberID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a group member")
	}
	return s.repo.AddGroupMember(groupID, ownerID, memberID)
}

func (s *GroupService) DeleteGroup(groupID, ownerID int) error {
	return s.repo.DeleteGroup(groupID, ownerID)
}

func (s *GroupService) LeaveGroup(groupID, userID int) error {
	return s.repo.LeaveGroup(groupID, userID)
}

func (s *GroupService) ConfirmMemberTime(groupID, userID int, confirmedTime time.Time) error {
	return s.repo.ConfirmMemberTime(groupID, userID, confirmedTime)
}

func (s *GroupService) GetGroupFreeTime(groupID, userID int, date time.Time) ([]model.TimeSlot, error) {
	if _, err := s.repo.GetGroupDetails(groupID, userID); err != nil {
		return nil, err
	}
	memberIDs, err := s.repo.GetMemberUserIDs(groupID)
	if err != nil {
		return nil, err
	}

	eventSvc := &EventService{repo: s.eventRepo}
	return eventSvc.GetFreeTimeForUsers(memberIDs, date)
}
