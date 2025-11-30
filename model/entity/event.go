package entity

import "time"

type EventStatus string

const (
	StatusPending  EventStatus = "pending"
	StatusAccepted EventStatus = "accepted"
	StatusDeclined EventStatus = "declined"
)

func (s EventStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusAccepted, StatusDeclined:
		return true
	}
	return false
}

type Event struct {
	ID          string                 `json:"id"`
	InitiatorID string                 `json:"initiatorId"`
	TeamID      string                 `json:"teamId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"createdAt"`
	StartsAt    time.Time              `json:"startsAt"`
	Duration    int64                  `json:"duration"`
	Statuses    map[string]EventStatus `json:"statuses"`
}

func NewEvent(id, initiatorId, teamId, name, description string, startsAt time.Time, duration int64, teamMembers []string) *Event {
	statuses := make(map[string]EventStatus, len(teamMembers))
	for _, member := range teamMembers {
		statuses[member] = EventStatus(StatusPending)
	}
	return &Event{
		ID:          id,
		InitiatorID: initiatorId,
		TeamID:      teamId,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		StartsAt:    startsAt,
		Duration:    duration,
		Statuses:    statuses,
	}
}
