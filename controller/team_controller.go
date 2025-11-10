package controller

import (
	"net/http"
	"strconv"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

type TeamController struct {
	teamService TeamServiceInterface
}

func NewTeamController() *TeamController {
	return &TeamController{
		teamService: service.NewTeamService(),
	}
}

func NewTeamControllerWithService(teamService TeamServiceInterface) *TeamController {
	return &TeamController{
		teamService: teamService,
	}
}

type TeamServiceInterface interface {
	CreateTeam(request *dto.TeamRequest) (*dto.TeamResponse, error)
	AddUserToTeam(idUser string, idTeam string) error
	GetTeamById(id string) (*entity.Team, error)
	GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error)
	GetTeamsByName(name string) ([]*entity.Team, error)
	GetAll() ([]*entity.Team, error)
	Update(team *entity.Team) error
	Delete(id string) error
}

// NewTeam
//
//	@Summary		Create a new team
//	@Description	Create a new team with the provided details
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.TeamRequest	true	"Team details"
//	@Success		201		{object}	dto.TeamResponse
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams [post]
func (tc *TeamController) NewTeam(c *gin.Context) {
	var request dto.TeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := tc.teamService.CreateTeam(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetTeam
//
//	@Summary		Get a team by ID
//	@Description	Get team details by ID
//	@Produce		json
//	@Param			id	path		string	true	"Team ID"
//	@Success		200	{object}	entity.Team
//	@Failure		404	{object}	map[string]interface{}	"Team not found"
//	@Router			/teams/{id} [get]
func (tc *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")
	team, err := tc.teamService.GetTeamById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetAllTeams
//
//	@Summary		Get all teams
//	@Description	Get a list of all teams
//	@Produce		json
//	@Success		200	{array}		entity.Team
//	@Failure		500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams [get]
func (tc *TeamController) GetAllTeams(c *gin.Context) {
	teams, err := tc.teamService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// AddUserToTeam
//
//	@Summary		Add a user to a team
//	@Description	Add a user to a team by providing user ID and team ID
//	@Accept			json
//	@Produce		json
//	@Param			request	body		map[string]string		true	"User ID and Team ID"
//	@Success		200		{object}	map[string]interface{}	"User added to team"
//	@Failure		400		{object}	map[string]interface{}	"Bad Request: Invalid request body or missing userId or teamId"
//	@Router			/teams/addUserToTeam [post]
func (tc *TeamController) AddUserToTeam(c *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
		TeamID string `json:"teamId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.UserID == "" || req.TeamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId or teamId"})
		return
	}

	if err := tc.teamService.AddUserToTeam(req.UserID, req.TeamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to team"})
}

// GetXTeamsByPrefix
//
//	@Summary		Get X teams by prefix
//	@Description	Get a list of X teams that start with the specified prefix
//	@Produce		json
//	@Param			prefix	query		string	true	"Prefix to search for"
//	@Param			limit	query		int		true	"Number of teams to retrieve"
//	@Success		200		{array}		entity.Team
//	@Failure		400		{object}	map[string]interface{}	"Bad Request: Missing prefix or limit query parameters, or limit is NaN"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams/search [get]
func (tc *TeamController) GetXTeamsByPrefix(c *gin.Context) {
	prefix := c.Query("prefix")
	xStr := c.Query("limit")

	if prefix == "" || xStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing prefix or limit query parameters"})
		return
	}

	x, err := strconv.Atoi(xStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be a number"})
		return
	}

	teams, err := tc.teamService.GetXTeamsByPrefix(prefix, x)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeamsByName
//
//	@Summary		Get teams by name
//	@Description	Get a list of teams that match the specified name
//	@Produce		json
//	@Param			name	query		string	true	"Name to search for"
//	@Success		200		{array}		entity.Team
//	@Failure		400		{object}	map[string]interface{}	"Bad Request: Missing 'name' query parameter"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams/by-name [get]
func (tc *TeamController) GetTeamsByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'name' query parameter"})
		return
	}

	teams, err := tc.teamService.GetTeamsByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// UpdateTeam
//
//	@Summary		Update a team
//	@Description	Update team details by providing team ID and updated details
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Team ID"
//	@Param			team	body		entity.Team	true	"Updated team details"
//	@Success		200		{object}	entity.Team
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams/{id} [put]
func (tc *TeamController) UpdateTeam(c *gin.Context) {
	var team entity.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.teamService.Update(&team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// DeleteTeam
//
//	@Summary		Delete a team
//	@Description	Delete a team by providing team ID
//	@Produce		json
//	@Param			id	path		string					true	"Team ID"
//	@Success		200	{object}	map[string]interface{}	"Team deleted"
//	@Failure		400	{object}	map[string]interface{}	"Bad Request: Missing team ID"
//	@Failure		500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams/{id} [delete]
func (tc *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing team ID"})
		return
	}

	if err := tc.teamService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted"})
}
