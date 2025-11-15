package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService service.MessageServiceInterface
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService: service.NewMessageService(),
	}
}

func NewMessageControllerWithService(messageService service.MessageServiceInterface) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
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
}

func (mc *MessageController) EditMessage(c *gin.Context) {
	// TODO: implement
}

func (mc *MessageController) DeleteMessage(c *gin.Context) {
	// TODO: implement
}
