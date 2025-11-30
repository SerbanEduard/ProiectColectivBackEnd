package service

import (
	"fmt"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
)

const (
	InvalidStatus = "invalid status"
)

type EventServiceInterface interface {
	CreateEvent(request *dto.CreateEventRequest) (*dto.EventDTO, error)
	GetEventById(id string) (*dto.EventDTO, error)
	GetEventsByTeamId(teamId string) ([]*dto.EventDTO, error)
	UpdateEventDetails(id string, request *dto.UpdateEventRequest) (*dto.EventDTO, error)
	UpdateUserStatus(id string, request *dto.UpdateEventStatusRequest) (*dto.EventDTO, error)
	DeleteEvent(id string) error
}

type EventService struct {
	userRepo  UserRepositoryInterface
	teamRepo  TeamRepositoryInterface
	eventRepo persistence.EventRepositoryInterface
}

func NewEventService() *EventService {
	return &EventService{
		userRepo:  persistence.NewUserRepository(),
		teamRepo:  persistence.NewTeamRepository(),
		eventRepo: persistence.NewEventRepository(),
	}
}

func (es *EventService) CreateEvent(req *dto.CreateEventRequest) (*dto.EventDTO, error) {
	if _, err := es.userRepo.GetByID(req.InitiatorID); err != nil {
		return nil, err
	}
	team, err := es.teamRepo.GetTeamById(req.TeamID)
	if err != nil {
		return nil, err
	}

	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	event := *entity.NewEvent(
		id,
		req.InitiatorID,
		req.TeamID,
		req.Name,
		req.Description,
		startsAt,
		req.Duration,
		team.UsersIds,
	)
	if err := es.eventRepo.Create(&event); err != nil {
		return nil, err
	}

	eventDTO := dto.NewEventDTO(&event)
	return eventDTO, nil
}

func (es *EventService) GetEventById(id string) (*dto.EventDTO, error) {
	event, err := es.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	eventDTO := dto.NewEventDTO(event)
	return eventDTO, nil
}

func (es *EventService) GetEventsByTeamId(teamId string) ([]*dto.EventDTO, error) {
	events, err := es.eventRepo.GetByTeamID(teamId)
	if err != nil {
		return nil, err
	}

	eventsDTO := make([]*dto.EventDTO, len(events))
	for i, event := range events {
		eventsDTO[i] = dto.NewEventDTO(event)
	}

	return eventsDTO, nil
}

func (es *EventService) UpdateEventDetails(id string, req *dto.UpdateEventRequest) (*dto.EventDTO, error) {
	event, err := es.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = req.Name
		event.Name = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
		event.Description = req.Description
	}
	if req.StartsAt != "" {
		if startsAt, err := time.Parse(time.RFC3339, req.StartsAt); err != nil {
			return nil, err
		} else {
			event.StartsAt = startsAt
		}
		updates["starts_at"] = req.StartsAt
	}
	if req.Duration != 0 {
		event.Duration = req.Duration
		updates["duration"] = req.Duration
	}

	if err := es.eventRepo.Update(id, updates); err != nil {
		return nil, err
	}

	return dto.NewEventDTO(event), nil
}

func (es *EventService) UpdateUserStatus(id string, req *dto.UpdateEventStatusRequest) (*dto.EventDTO, error) {
	event, err := es.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if _, err := es.userRepo.GetByID(req.UserID); err != nil {
		return nil, err
	}

	status := entity.EventStatus(req.Status)
	if !status.IsValid() {
		return nil, fmt.Errorf(InvalidStatus)
	}

	event.Statuses[req.UserID] = status
	updates := map[string]interface{}{
		"statuses": event.Statuses,
	}

	if err := es.eventRepo.Update(id, updates); err != nil {
		return nil, err
	}

	return dto.NewEventDTO(event), nil
}

func (es *EventService) DeleteEvent(id string) error {
	if _, err := es.eventRepo.GetByID(id); err != nil {
		return err
	}

	return es.eventRepo.Delete(id)
}
