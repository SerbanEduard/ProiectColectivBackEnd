package controller

import (
	"net/http"
	"slices"

	"github.com/SerbanEduard/ProiectColectivBackEnd/hub"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/utils"
	"github.com/gin-gonic/gin"
)

const (
	BadMessageTypeError  = "message type must be direct or team"
	MessageNotFoundError = "message not found"
	MissingParameter     = "missing parameter(s)"
)

type MessageRequestUnion struct {
	Direct *dto.DirectMessageRequest `json:"direct,omitempty"`
	Team   *dto.TeamMessageRequest   `json:"team,omitempty"`
}

type MessageController struct {
	messageService service.MessageServiceInterface
	teamService    TeamServiceInterface
	hub            *hub.Hub[hub.Message]
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService: service.NewMessageService(),
		teamService:    service.NewTeamService(),
		hub:            hub.NewHub[hub.Message](),
	}
}

func NewMessageControllerWithService(messageService service.MessageServiceInterface) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// Connect
//
//	@Summary	Connect the user to the message WebSocket
//	@Tags		messages
//	@Security	Bearer
//	@Success	101	{string}	string					"Switching Protocols - WebSocket connection established"
//	@Failure	400	{object}	map[string]interface{}	"Bad Request"
//	@Failure	401	{object}	map[string]string		"Unauthorized"
//	@Failure	500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/messages/connect [get]
func (mc *MessageController) Connect(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		// This should never happen if the auth middleware is used
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conn, err := hub.AcceptConnection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := hub.NewClient[hub.Message](userID, conn)
	mc.hub.Register(client)
}

// NewMessage
//
//	@Summary		Create and send a message
//	@Description	Create and send a message either to another user or to a team
//	@Tags			messages
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			type	query		string				true	"Message type (direct/team)"
//	@Param			request	body		MessageRequestUnion	true	"The message request (this is only for documentation purposes, the actual request should be either DirectMessageRequest or TeamMessageRequest)"
//	@Success		201		{object}	dto.MessageDTO
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/messages [post]
func (mc *MessageController) NewMessage(c *gin.Context) {
	message_type := c.Query("type")

	switch message_type {
	case "direct":
		var request dto.DirectMessageRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := mc.messageService.CreateDirectMessage(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)

		mc.hub.Send(request.ReceiverID, *hub.NewMessage(hub.DirectMessage, request))

	case "team":
		var request dto.TeamMessageRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := mc.messageService.CreateTeamMessage(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resp)

		// Send to team members via WebSocket
		team, _ := mc.teamService.GetTeamById(request.TeamId)
		userIDs := slices.DeleteFunc(team.UsersIds, func(uid string) bool {
			return uid != request.SenderID // Don't send to sender
		})
		mc.hub.SendMany(userIDs, *hub.NewMessage(hub.TeamBroadcast, request))

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": BadMessageTypeError})
	}
}

// GetMessage
//
//	@Summary	Get a message by ID
//	@Tags		messages
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The message ID"
//	@Success	200	{object}	dto.MessageDTO
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/messages/{id} [get]
func (mc *MessageController) GetMessage(c *gin.Context) {
	id := c.Param("id")
	message, err := mc.messageService.GetMessageByID(id)
	if err.Error() == entity.BadConversationKey {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": MessageNotFoundError})
		return
	}

	c.JSON(http.StatusOK, message)
}

// GetMessages
//
//	@Summary		Get all messages
//	@Description	Get messages between 2 users or within a team
//	@Tags			messages
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			type	query		string	true	"Messages type (direct/team)"
//	@Param			user1Id	query		string	false	"User1 ID (direct message)"
//	@Param			user2Id	query		string	false	"User2 ID (direct message)"
//	@Param			teamId	query		string	false	"Team ID (team message)"
//	@Success		200		{array}		dto.MessageDTO
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/messages [get]
func (mc *MessageController) GetMessages(c *gin.Context) {
	message_type := c.Query("type")

	switch message_type {
	case "direct":
		user1Id := c.Query("user1Id")
		user2Id := c.Query("user2Id")

		if user1Id == "" || user2Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": MissingParameter})
			return
		}

		resp, err := mc.messageService.GetDirectMessages(user1Id, user2Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	case "team":
		teamId := c.Query("teamId")
		if teamId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": MissingParameter})
			return
		}

		resp, err := mc.messageService.GetTeamMessages(teamId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": BadMessageTypeError})
	}
}

func (mc *MessageController) EditMessage(c *gin.Context) {
	// TODO: implement
}

func (mc *MessageController) DeleteMessage(c *gin.Context) {
	// TODO: implement
}
