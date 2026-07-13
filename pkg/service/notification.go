package service

import (
	"druna_server/pkg/repository"
	"encoding/json"
	"time"
)

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) EnqueueFriendRequest(targetUserID int, fromUsername string) {
	payload, _ := json.Marshal(map[string]string{
		"fromUsername": fromUsername,
	})
	_ = s.repo.Enqueue(targetUserID, "friend_request", string(payload))
}

func (s *NotificationService) EnqueueGroupConfirm(targetUserID int, groupID int) {
	payload, _ := json.Marshal(map[string]int{
		"groupId": groupID,
	})
	_ = s.repo.Enqueue(targetUserID, "group_confirm", string(payload))
}

func (s *NotificationService) EnqueueGroupEventCreated(targetUserID, groupID, eventID int, title string, startTime time.Time) {
	payload, _ := json.Marshal(map[string]interface{}{
		"groupId":   groupID,
		"eventId":   eventID,
		"title":     title,
		"startTime": startTime,
	})
	_ = s.repo.Enqueue(targetUserID, "group_event_created", string(payload))
}
