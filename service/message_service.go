package service

import (
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

type MessageService struct {
	userRepo    UserRepositoryInterface
	teamRepo    TeamRepositoryInterface
	messageRepo persistence.MessageRepositoryInterface
}

func NewMessageService() *MessageService {
	return &MessageService{
		userRepo:    persistence.NewUserRepository(),
		teamRepo:    persistence.NewTeamRepository(),
		messageRepo: persistence.NewMessageRepository(),
	}
}

func NewMessageServiceWithRepo(userRepo UserRepositoryInterface, teamRepo TeamRepositoryInterface, messageRepo persistence.MessageRepositoryInterface) *MessageService {
	return &MessageService{
		userRepo:    userRepo,
		teamRepo:    teamRepo,
		messageRepo: messageRepo,
	}
}

type MessageServiceInterface interface {
	CreateDirectMessage(request *dto.DirectMessageRequest) (*dto.MessageDTO, error)
	CreateTeamMessage(request *dto.TeamMessageRequest) (*dto.MessageDTO, error)
	GetMessageByID(id string) (*dto.MessageDTO, error)
	GetDirectMessages(user1Id, user2Id string) ([]*dto.MessageDTO, error)
	GetTeamMessages(teamId string) ([]*dto.MessageDTO, error)
}

func (ms *MessageService) CreateDirectMessage(request *dto.DirectMessageRequest) (*dto.MessageDTO, error) {
	if err := validator.ValidateDirectMessageRequest(request); err != nil {
		return nil, err
	}

	sender, err := ms.userRepo.GetByID(request.SenderID)
	if err != nil {
		return nil, fmt.Errorf("sender not found")
	}

	if _, err := ms.userRepo.GetByID(request.ReceiverID); err != nil {
		return nil, fmt.Errorf("receiver not found")
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	message := *entity.NewMessage(
		id,
		request.SenderID,
		entity.GetConversationKey(request.SenderID, request.ReceiverID),
		"",
		request.TextContent,
	)
	if err := ms.messageRepo.Create(&message); err != nil {
		return nil, err
	}

	senderDTO := dto.NewSenderDTO(sender)
	dtoMessage := dto.NewMessageDTO(message.ID, request.ReceiverID, "", message.TextContent, message.SentAt, *senderDTO)
	return dtoMessage, nil
}

func (ms *MessageService) CreateTeamMessage(request *dto.TeamMessageRequest) (*dto.MessageDTO, error) {
	if err := validator.ValidateTeamMessageRequest(request); err != nil {
		return nil, err
	}

	sender, err := ms.userRepo.GetByID(request.SenderID)
	if err != nil {
		return nil, fmt.Errorf("sender not found")
	}

	if _, err := ms.teamRepo.GetTeamById(request.TeamId); err != nil {
		return nil, fmt.Errorf("team not found")
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	message := *entity.NewMessage(
		id,
		request.SenderID,
		"",
		request.TeamId,
		request.TextContent,
	)
	if err := ms.messageRepo.Create(&message); err != nil {
		return nil, err
	}

	senderDTO := dto.NewSenderDTO(sender)
	dtoMessage := dto.NewMessageDTO(message.ID, "", request.TeamId, message.TextContent, message.SentAt, *senderDTO)
	return dtoMessage, nil
}

func (ms *MessageService) GetMessageByID(id string) (*dto.MessageDTO, error) {
	message, err := ms.messageRepo.GetByID(id)
	receiverId, key_err := entity.GetReceiverIdFromKey(message.SenderID, message.ConversationKey)
	if message.ConversationKey != "" && key_err != nil {
		return nil, err
	}

	sender, err := ms.userRepo.GetByID(message.SenderID)
	if err != nil {
		return nil, fmt.Errorf("sender not found")
	}

	senderDTO := dto.NewSenderDTO(sender)
	dtoMessage := dto.NewMessageDTO(message.ID, receiverId, message.TeamID, message.TextContent, message.SentAt, *senderDTO)
	return dtoMessage, err
}

func (ms *MessageService) GetDirectMessages(user1Id, user2Id string) ([]*dto.MessageDTO, error) {
	if _, err := ms.userRepo.GetByID(user1Id); err != nil {
		return nil, fmt.Errorf("user1 not found")
	}
	if _, err := ms.userRepo.GetByID(user2Id); err != nil {
		return nil, fmt.Errorf("user2 not found")
	}

	messages, err := ms.messageRepo.GetByConversation(user1Id, user2Id)
	dtoMessages := []*dto.MessageDTO{}
	for _, message := range messages {
		receiverId, key_err := entity.GetReceiverIdFromKey(message.SenderID, message.ConversationKey)
		if message.ConversationKey != "" && key_err != nil {
			return nil, err
		}

		sender, err := ms.userRepo.GetByID(message.SenderID)
		if err != nil {
			return nil, fmt.Errorf("sender not found")
		}

		senderDTO := dto.NewSenderDTO(sender)
		dtoMessage := dto.NewMessageDTO(message.ID, receiverId, message.TeamID, message.TextContent, message.SentAt, *senderDTO)
		dtoMessages = append(dtoMessages, dtoMessage)
	}
	return dtoMessages, err
}

func (ms *MessageService) GetTeamMessages(teamId string) ([]*dto.MessageDTO, error) {
	if _, err := ms.teamRepo.GetTeamById(teamId); err != nil {
		return nil, fmt.Errorf("team not found")
	}

	messages, err := ms.messageRepo.GetByTeamID(teamId)
	dtoMessages := []*dto.MessageDTO{}
	for _, message := range messages {

		sender, err := ms.userRepo.GetByID(message.SenderID)
		if err != nil {
			return nil, fmt.Errorf("sender not found")
		}

		senderDTO := dto.NewSenderDTO(sender)
		dtoMessage := dto.NewMessageDTO(message.ID, "", message.TeamID, message.TextContent, message.SentAt, *senderDTO)
		dtoMessages = append(dtoMessages, dtoMessage)
	}
	return dtoMessages, err
}
