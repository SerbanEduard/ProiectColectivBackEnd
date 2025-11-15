package persistence

import (
	"context"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const (
	messagesCollection = "messages"
	convKeyField       = "convKey"
)

type MessageRepositoryInterface interface {
	Create(message *entity.Message) error
	GetByConversation(user1Id, user2Id string) ([]*entity.Message, error)
	GetByTeam(teamId string) ([]*entity.Message, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
}

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (mr *MessageRepository) Create(message *entity.Message) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(messagesCollection + "/" + message.ID)
	return ref.Set(ctx, message)
}

func (mr *MessageRepository) GetByConversation(user1Id, user2Id string) ([]*entity.Message, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(messagesCollection)

	query := ref.OrderByChild(convKeyField).EqualTo(entity.GetConversationKey(user1Id, user2Id))
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, 0, len(results))
	for _, r := range results {
		var message entity.Message
		if err := r.Unmarshal(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (mr *MessageRepository) GetByTeam(teamId string) ([]*entity.Message, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(messagesCollection)

	query := ref.OrderByChild("teamId").EqualTo(teamId)
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, 0, len(results))
	for _, r := range results {
		var message entity.Message
		if err := r.Unmarshal(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (mr *MessageRepository) Update(id string, updates map[string]interface{}) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(messagesCollection + "/" + id)
	return ref.Update(ctx, updates)
}

func (mr *MessageRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(messagesCollection + "/" + id)
	return ref.Delete(ctx)
}
