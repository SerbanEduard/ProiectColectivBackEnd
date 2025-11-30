package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

const (
	EventDeleted = "Event deleted succesfully"
)

type EventController struct {
	eventService service.EventServiceInterface
	teamService  TeamServiceInterface
}

func NewEventController() *EventController {
	return &EventController{
		eventService: service.NewEventService(),
	}
}

// NewEvent
//
//	@Summary	Create new event
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.CreateEventRequest	true	"Create event request"
//	@Success	201		{object}	dto.EventDTO
//	@Failure	400		{object}	map[string]interface{}	"Bad Request"
//	@Failure	500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/events [post]
func (ec *EventController) NewEvent(c *gin.Context) {
	var request dto.CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ec.eventService.CreateEvent(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetEvent
//
//	@Summary	Get event by id
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Event ID"
//	@Success	200	{object}	dto.EventDTO
//	@Failure	400	{object}	map[string]interface{}	"Bad Request"
//	@Failure	500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/events/{id} [get]
func (ec *EventController) GetEvent(c *gin.Context) {
	id := c.Param("id")
	event, err := ec.eventService.GetEventById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetEvents
//
//	@Summary	Get events by team id
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		teamId	query		string	true	"Team ID"
//	@Success	200		{object}	[]dto.EventDTO
//	@Failure	400		{object}	map[string]interface{}	"Bad Request"
//	@Failure	500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/events [get]
func (ec *EventController) GetEvents(c *gin.Context) {
	teamId := c.Query("teamId")

	if teamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": MissingParameter})
		return
	}
	if _, err := ec.teamService.GetTeamById(teamId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	events, err := ec.eventService.GetEventsByTeamId(teamId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEventDetails
//
//	@Summary		Update event details
//	@Description	Update event name, description, start time and/or duration
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Event ID"
//	@Param			request	body		dto.UpdateEventRequest	true	"Update event request"
//	@Success		200		{object}	dto.EventDTO
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/events/{id} [patch]
func (ec *EventController) UpdateEventDetails(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ec.eventService.UpdateEventDetails(id, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserStatus
//
//	@Summary	Update user status for event
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string							true	"Event ID"
//	@Param		request	body		dto.UpdateEventStatusRequest	true	"Update event status request"
//	@Success	200		{object}	dto.EventDTO
//	@Failure	400		{object}	map[string]interface{}	"Bad Request"
//	@Failure	500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/events/{id}/status [patch]
func (ec *EventController) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateEventStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := ec.eventService.UpdateUserStatus(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent
//
//	@Summary	Delete event by id
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string					true	"Event ID"
//	@Success	200	{object}	map[string]interface{}	"Event deleted successfully"}]
//	@Failure	500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router		/events/{id} [delete]
func (ec *EventController) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	if err := ec.eventService.DeleteEvent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": EventDeleted})
}
