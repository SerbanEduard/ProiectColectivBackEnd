package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

type TeamRequestController struct {
	teamRequestService *service.TeamRequestService
}

func NewTeamRequestController() *TeamRequestController {
	return &TeamRequestController{
		teamRequestService: service.NewTeamRequestService(),
	}
}

// CreateTeamRequest
//
//	@Summary		Create a new team request
//	@Description	Creates a join request for a user to join a team.
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.TeamRequestCreateDTO	true	"Team and User IDs"
//	@Success		201		{object}	dto.TeamRequestItemDTO
//	@Failure		400		{object}	map[string]string	"Invalid request or validation error"
//	@Failure		500		{object}	map[string]string	"Internal Server Error"
//	@Router			/teamRequests [post]
func (tc *TeamRequestController) CreateTeamRequest(c *gin.Context) {
	var req dto.TeamRequestCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := tc.teamRequestService.CreateTeamRequest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.NewTeamRequestItemDTO(created))
}

// AcceptTeamRequest
//
//	@Summary		Accept a team request
//	@Description	Accepts a pending team request, adds the user to the team, and deletes the request.
//	@Security		Bearer
//	@Produce		json
//	@Param			id	path		string	true	"Team Request ID"
//	@Success		200	{object}	dto.AddUserToTeamResponse
//	@Failure		400	{object}	map[string]string	"Bad Request or Not Found"
//	@Failure		500	{object}	map[string]string	"Internal Server Error"
//	@Router			/teamRequests/{id}/accept [put]
func (tc *TeamRequestController) AcceptTeamRequest(c *gin.Context) {
	id := c.Param("id")
	user, team, err := tc.teamRequestService.AcceptTeamRequest(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.NewAddUserToTeamResponse(*user, *team))
}

// RejectTeamRequest
//
//	@Summary		Reject a team request
//	@Description	Rejects and deletes a pending team request.
//	@Security		Bearer
//	@Produce		json
//	@Param			id	path		string				true	"Team Request ID"
//	@Success		200	{object}	map[string]string	"Team request rejected"
//	@Failure		400	{object}	map[string]string	"Bad Request or Not Found"
//	@Router			/teamRequests/{id}/reject [delete]
func (tc *TeamRequestController) RejectTeamRequest(c *gin.Context) {
	id := c.Param("id")
	if err := tc.teamRequestService.RejectTeamRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team request rejected"})
}

// GetAllTeamRequests
//
//	@Summary		Get all team requests
//	@Description	Fetches all pending team requests in the system.
//	@Security		Bearer
//	@Produce		json
//	@Success		200	{object}	dto.TeamRequestsResponseDTO
//	@Failure		500	{object}	map[string]string	"Internal Server Error"
//	@Router			/teamRequests [get]
func (tc *TeamRequestController) GetAllTeamRequests(c *gin.Context) {
	reqs, err := tc.teamRequestService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.NewTeamRequestsResponseDTO(reqs))
}

// GetTeamRequestsByUser
//
//	@Summary		Get all team requests for a specific user
//	@Description	Fetches all pending team requests created by a given user.
//	@Security		Bearer
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	dto.TeamRequestsResponseDTO
//	@Failure		400		{object}	map[string]string	"Bad Request"
//	@Failure		500		{object}	map[string]string	"Internal Server Error"
//	@Router			/teamRequests/user/{userId} [get]
func (tc *TeamRequestController) GetTeamRequestsByUser(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
		return
	}

	reqs, err := tc.teamRequestService.GetByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.NewTeamRequestsResponseDTO(reqs))
}
