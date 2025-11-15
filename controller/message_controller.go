package controller

import (
	"net/http"
	"slices"

	"github.com/SerbanEduard/ProiectColectivBackEnd/hub"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/utils"
	"github.com/gin-gonic/gin"
)

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

// NewDirectMessage
//
//	@Summary	Create and send a direct message
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.DirectMessagesRequest	true	"The direct message request"
//	@Success	201		{object}	entity.Message
//	@Failure	400		{object}	map[string]interface{}	"Bad Request"
//	@Failure	500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/messages/direct [post]
func (mc *MessageController) NewDirectMessage(c *gin.Context) {
	var request dto.DirectMessagesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := mc.messageService.CreateDirectMessage(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)

	// Send to receiver via WebSocket
	mc.hub.Send(request.ReceiverId, *hub.NewMessage(hub.DirectMessage, request))
}

// NewTeamMessage
//
//	@Summary	Create and send a team message
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.TeamMessagesRequest	true	"The team message request"
//	@Success	201		{object}	entity.Message
//	@Failure	400		{object}	map[string]interface{}	"Bad Request"
//	@Failure	500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/messages/team [post]
func (mc *MessageController) NewTeamMessage(c *gin.Context) {
	var request dto.TeamMessagesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := mc.messageService.CreateTeamMessage(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)

	// Send to team members via WebSocket
	team, err := mc.teamService.GetTeamById(request.TeamId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	userIDs := slices.DeleteFunc(team.UsersIds, func(uid string) bool {
		return uid != request.SenderId // Don't send to sender
	})
	mc.hub.SendMany(userIDs, *hub.NewMessage(hub.TeamBroadcast, request))
}

func (mc *MessageController) EditMessage(c *gin.Context) {
	// TODO: implement
}

func (mc *MessageController) DeleteMessage(c *gin.Context) {
	// TODO: implement
}
