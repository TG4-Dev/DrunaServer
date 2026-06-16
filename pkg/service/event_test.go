package service

import (
	"druna_server/pkg/model"
	"testing"
	"time"
)

type mockEventRepo struct {
	events []model.Event
}

func (m *mockEventRepo) CreateEvent(event model.Event) (int, error) { return 1, nil }
func (m *mockEventRepo) DeleteEvent(userID, eventID int) error      { return nil }
func (m *mockEventRepo) GetEventList(userID int) ([]model.Event, error) {
	return m.events, nil
}

func TestGetFreeTimeNoEvents(t *testing.T) {
	svc := NewEventService(&mockEventRepo{})
	date := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)

	slots, err := svc.GetFreeTime(1, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slots) != 1 {
		t.Fatalf("expected 1 free slot, got %d", len(slots))
	}
	if !slots[0].Start.Equal(date) {
		t.Fatalf("unexpected start: %v", slots[0].Start)
	}
}

func TestGetFreeTimeWithBusySlot(t *testing.T) {
	date := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	repo := &mockEventRepo{
		events: []model.Event{
			{
				StartTime: date.Add(10 * time.Hour),
				EndTime:   date.Add(12 * time.Hour),
			},
		},
	}
	svc := NewEventService(repo)

	slots, err := svc.GetFreeTime(1, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slots) != 2 {
		t.Fatalf("expected 2 free slots, got %d", len(slots))
	}
	if !slots[0].End.Equal(date.Add(10 * time.Hour)) {
		t.Fatalf("unexpected first slot end: %v", slots[0].End)
	}
	if !slots[1].Start.Equal(date.Add(12 * time.Hour)) {
		t.Fatalf("unexpected second slot start: %v", slots[1].Start)
	}
}

func TestCreateEventValidation(t *testing.T) {
	svc := NewEventService(&mockEventRepo{})
	start := time.Now()
	_, err := svc.CreateEvent(model.Event{
		StartTime: start,
		EndTime:   start.Add(-time.Hour),
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
