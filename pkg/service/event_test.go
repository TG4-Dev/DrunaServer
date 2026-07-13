package service

import (
	"druna_server/pkg/model"
	"testing"
	"time"
)

type mockEventRepo struct {
	events []model.Event
}

func (m *mockEventRepo) CreateEvent(event model.Event) (int, error)      { return 1, nil }
func (m *mockEventRepo) UpdateEvent(userID int, event model.Event) error { return nil }
func (m *mockEventRepo) DeleteEvent(userID, eventID int) error           { return nil }
func (m *mockEventRepo) HasOverlappingEvent(userID int, start, end time.Time, excludeID int) (bool, error) {
	return false, nil
}
func (m *mockEventRepo) GetEventList(userID int) ([]model.Event, error) { return m.events, nil }
func (m *mockEventRepo) GetEventListFiltered(userID int, filter model.EventFilter) ([]model.Event, error) {
	return m.events, nil
}
func (m *mockEventRepo) CountEvents(userID int, filter model.EventFilter) (int, error) {
	return len(m.events), nil
}
func (m *mockEventRepo) GetBusyEventsForUsers(userIDs []int, dateFrom, dateTo time.Time) (map[int][]model.Event, error) {
	return map[int][]model.Event{}, nil
}
func (m *mockEventRepo) CreateGroupEvent(event model.Event) (int, error) { return 1, nil }
func (m *mockEventRepo) UpdateGroupEvent(groupID, eventID int, event model.Event) error {
	return nil
}
func (m *mockEventRepo) DeleteGroupEvent(groupID, eventID int) error { return nil }
func (m *mockEventRepo) GetGroupEventByID(groupID, eventID int) (model.Event, error) {
	return model.Event{}, nil
}
func (m *mockEventRepo) GetGroupEvents(groupID int, filter model.EventFilter) ([]model.Event, error) {
	return m.events, nil
}
func (m *mockEventRepo) CountGroupEvents(groupID int, filter model.EventFilter) (int, error) {
	return len(m.events), nil
}
func (m *mockEventRepo) HasOverlappingGroupEvent(groupID int, start, end time.Time, excludeID int) (bool, error) {
	return false, nil
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

func TestIntersectTimeSlots(t *testing.T) {
	day := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	a := []model.TimeSlot{{Start: day.Add(9 * time.Hour), End: day.Add(12 * time.Hour)}}
	b := []model.TimeSlot{{Start: day.Add(10 * time.Hour), End: day.Add(14 * time.Hour)}}
	result := IntersectTimeSlots(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 intersection, got %d", len(result))
	}
	if !result[0].Start.Equal(day.Add(10 * time.Hour)) {
		t.Fatalf("unexpected intersection start")
	}
}
