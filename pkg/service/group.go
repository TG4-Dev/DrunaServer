package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
)

type GroupService struct {
	repo repository.Group
}

func NewGroupService(repo repository.Group) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(input model.Group) (int, error) {
	return s.repo.CreateGroup(input)
}
