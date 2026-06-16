package service

import (
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
	"sort"
	"time"
)

type EventService struct {
	repo repository.Event
}

func NewEventService(repo repository.Event) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(event model.Event) (int, error) {
	if !event.EndTime.After(event.StartTime) {
		return 0, errors.New("end time must be after start time")
	}
	return s.repo.CreateEvent(event)
}

func (s *EventService) DeleteEvent(userID, eventID int) error {
	return s.repo.DeleteEvent(userID, eventID)
}

func (s *EventService) GetEventList(userID int) ([]model.Event, error) {
	return s.repo.GetEventList(userID)
}

func (s *EventService) GetFreeTime(userID int, date time.Time) ([]model.TimeSlot, error) {
	events, err := s.repo.GetEventList(userID)
	if err != nil {
		return nil, err
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var busy []model.TimeSlot
	for _, event := range events {
		if event.EndTime.Before(dayStart) || !event.StartTime.Before(dayEnd) {
			continue
		}

		start := event.StartTime
		if start.Before(dayStart) {
			start = dayStart
		}

		end := event.EndTime
		if end.After(dayEnd) {
			end = dayEnd
		}

		if start.Before(end) {
			busy = append(busy, model.TimeSlot{Start: start, End: end})
		}
	}

	sort.Slice(busy, func(i, j int) bool {
		return busy[i].Start.Before(busy[j].Start)
	})

	merged := make([]model.TimeSlot, 0, len(busy))
	for _, slot := range busy {
		if len(merged) == 0 {
			merged = append(merged, slot)
			continue
		}
		last := &merged[len(merged)-1]
		if !slot.Start.After(last.End) {
			if slot.End.After(last.End) {
				last.End = slot.End
			}
			continue
		}
		merged = append(merged, slot)
	}

	free := make([]model.TimeSlot, 0)
	cursor := dayStart
	for _, slot := range merged {
		if cursor.Before(slot.Start) {
			free = append(free, model.TimeSlot{Start: cursor, End: slot.Start})
		}
		if slot.End.After(cursor) {
			cursor = slot.End
		}
	}
	if cursor.Before(dayEnd) {
		free = append(free, model.TimeSlot{Start: cursor, End: dayEnd})
	}

	return free, nil
}
