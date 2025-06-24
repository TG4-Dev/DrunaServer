package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
)

type EventService struct {
	repo repository.Event
}

func NewEventService(repo repository.Event) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event model.Event) (int, error) {
	return s.repo.CreateEvent(event)
}

func (s *EventService) DeleteEvent(userID, eventID int) error {
	return s.repo.DeleteEvent(userID, eventID)
}

func (s *EventService) GetEventList(userID int) ([]model.Event, error) {
	return s.repo.GetEventList(userID)
}
