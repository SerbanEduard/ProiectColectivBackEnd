package dto

import (
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

type CreateEventRequest struct {
	InitiatorID string `json:"initiatorId"`
	TeamID      string `json:"teamId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartsAt    string `json:"startsAt"`
	Duration    int64  `json:"duration"`
}

func NewCreateEventRequest(initiatorID, teamID, name, description, startsAt string, duration int64) *CreateEventRequest {
	return &CreateEventRequest{
		InitiatorID: initiatorID,
		TeamID:      teamID,
		Name:        name,
		Description: description,
		StartsAt:    startsAt,
		Duration:    duration,
	}
}

type UpdateEventRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartsAt    string `json:"startsAt"`
	Duration    int64  `json:"duration"`
}

func NewUpdateEventRequest(name, description, startsAt string, duration int64) *UpdateEventRequest {
	return &UpdateEventRequest{
		Name:        name,
		Description: description,
		StartsAt:    startsAt,
		Duration:    duration,
	}
}

type UpdateEventStatusRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

func NewUpdateEventStatusRequest(userID, status string) *UpdateEventStatusRequest {
	return &UpdateEventStatusRequest{
		UserID: userID,
		Status: status,
	}
}

type EventDTO struct {
	ID            string `json:"id"`
	InitiatorID   string `json:"initiatorId"`
	TeamID        string `json:"teamId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	StartsAt      string `json:"startsAt"`
	Duration      int64  `json:"duration"`
	PendingCount  int64  `json:"pendingCount"`
	AcceptedCount int64  `json:"acceptedCount"`
	DeclinedCount int64  `json:"declinedCount"`
}

func NewEventDTO(event *entity.Event) *EventDTO {
	pendingCount, acceptedCount, declinedCount := GetStatusCount(event.Statuses)
	return &EventDTO{
		ID:            event.ID,
		InitiatorID:   event.InitiatorID,
		TeamID:        event.TeamID,
		Name:          event.Name,
		Description:   event.Description,
		StartsAt:      event.StartsAt.Format(time.RFC3339),
		Duration:      event.Duration,
		PendingCount:  pendingCount,
		AcceptedCount: acceptedCount,
		DeclinedCount: declinedCount,
	}
}

func GetStatusCount(statuses map[string]entity.EventStatus) (int64, int64, int64) {
	pendingCount := int64(0)
	acceptedCount := int64(0)
	declinedCount := int64(0)
	for _, status := range statuses {
		switch status {
		case entity.StatusPending:
			pendingCount++
		case entity.StatusAccepted:
			acceptedCount++
		case entity.StatusDeclined:
			declinedCount++
		}
	}
	return pendingCount, acceptedCount, declinedCount
}
