package persistence

import (
	"context"
	"errors"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const (
	eventsCollection = "events"
	EventNotFound    = "event not found"
)

type EventRepositoryInterface interface {
	Create(event *entity.Event) error
	GetByID(id string) (*entity.Event, error)
	GetByTeamID(teamId string) ([]*entity.Event, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
}

type EventRepository struct {
}

func NewEventRepository() EventRepositoryInterface {
	return &EventRepository{}
}

func (er *EventRepository) Create(event *entity.Event) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(eventsCollection + "/" + event.ID)
	return ref.Set(ctx, event)
}

func (er *EventRepository) GetByID(id string) (*entity.Event, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(eventsCollection + "/" + id)

	var event entity.Event
	if err := ref.Get(ctx, &event); err != nil {
		return nil, err
	}
	if event.ID == "" {
		return nil, errors.New(EventNotFound)
	}
	return &event, nil
}

func (er *EventRepository) GetByTeamID(teamId string) ([]*entity.Event, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(eventsCollection)

	query := ref.OrderByChild("teamId").EqualTo(teamId)
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]*entity.Event, 0, len(results))
	for _, r := range results {
		var event entity.Event
		if err := r.Unmarshal(&event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (er *EventRepository) Update(id string, updates map[string]interface{}) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(eventsCollection + "/" + id)
	return ref.Update(ctx, updates)
}

func (er *EventRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(eventsCollection + "/" + id)
	return ref.Delete(ctx)
}
