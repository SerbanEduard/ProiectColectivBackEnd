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
	CreateDirectMessage(request *dto.DirectMessagesRequest) (*entity.Message, error)
	CreateTeamMessage(request *dto.TeamMessagesRequest) (*entity.Message, error)
}

func (ms *MessageService) CreateDirectMessage(request *dto.DirectMessagesRequest) (*entity.Message, error) {
	if err := validator.ValidateDirectMessageRequest(request); err != nil {
		return nil, err
	}

	if _, err := ms.userRepo.GetByID(request.SenderId); err != nil {
		return nil, fmt.Errorf("sender not found")
	}

	if _, err := ms.userRepo.GetByID(request.ReceiverId); err != nil {
		return nil, fmt.Errorf("receiver not found")
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	message := *entity.NewMessage(
		id,
		request.SenderId,
		entity.GetConversationKey(request.SenderId, request.ReceiverId),
		"",
		request.TextContent,
	)
	if err := ms.messageRepo.Create(&message); err != nil {
		return nil, err
	}
	return &message, nil
}

func (ms *MessageService) CreateTeamMessage(request *dto.TeamMessagesRequest) (*entity.Message, error) {
	if err := validator.ValidateTeamMessageRequest(request); err != nil {
		return nil, err
	}

	if _, err := ms.userRepo.GetByID(request.SenderId); err != nil {
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
		request.SenderId,
		"",
		request.TeamId,
		request.TextContent,
	)
	if err := ms.messageRepo.Create(&message); err != nil {
		return nil, err
	}
	return &message, nil
}
