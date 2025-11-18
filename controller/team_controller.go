package controller

import (
	"net/http"
	"strconv"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

const (
	TeamNotFoundError           = "Team not found"
	LimitParameterRequiredError = "limit parameter required when using prefix"
	LimitMustBeANumberError     = "limit must be a number"
	InvalidRequestBodyError     = "Invalid request body"
	EmptyParametersError        = "Invalid request body, empty parameters"
	MissingTeamIDError          = "Missing team ID"
	UserAddedToTeamMessage      = "User added to team"
	UserDeletedFromTeamMessage  = "User deleted from team"
	TeamDeletedMessage          = "Team deleted"
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
	CreateTeam(request *dto.TeamRequest) (*entity.Team, error)
	AddUserToTeam(idUser string, idTeam string) (*entity.User, *entity.Team, error)
	DeleteUserFromTeam(idUser string, idTeam string) (*entity.User, *entity.Team, error)
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
//	@Security		Bearer
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.TeamRequest	true	"Team details"
//	@Success		201		{object}	entity.Team
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
//	@Security		Bearer
//	@Produce		json
//	@Param			id	path		string	true	"Team ID"
//	@Success		200	{object}	entity.Team
//	@Failure		404	{object}	map[string]interface{}	"Team not found"
//	@Router			/teams/{id} [get]
func (tc *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")
	team, err := tc.teamService.GetTeamById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": TeamNotFoundError})
		return
	}

	c.JSON(http.StatusOK, team)
}

// GetAllTeams
//
//	@Summary		Get teams with optional filtering
//	@Description	Get teams - all teams, by name, or by prefix with limit
//	@Security		Bearer
//	@Produce		json
//	@Param			name	query		string	false	"Filter by exact name"
//	@Param			prefix	query		string	false	"Filter by name prefix"
//	@Param			limit	query		int		false	"Limit results (required with prefix)"
//	@Success		200		{array}		entity.Team
//	@Failure		400		{object}	map[string]interface{}	"Bad Request"
//	@Failure		500		{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams [get]
func (tc *TeamController) GetAllTeams(c *gin.Context) {
	name := c.Query("name")
	prefix := c.Query("prefix")
	limitStr := c.Query("limit")

	// Filter by exact name
	if name != "" {
		teams, err := tc.teamService.GetTeamsByName(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, teams)
		return
	}

	// Filter by prefix with limit
	if prefix != "" {
		if limitStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": LimitParameterRequiredError})
			return
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": LimitMustBeANumberError})
			return
		}
		teams, err := tc.teamService.GetXTeamsByPrefix(prefix, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, teams)
		return
	}

	// Get all teams (no filters)
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
//	@Description	Adds a user to a team by providing user ID and team ID
//
//	@Security		Bearer
//
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserToTeamRequest	true	"User ID and Team ID"
//	@Success		200		{object}	dto.AddUserToTeamResponse
//	@Failure		400		{object}	map[string]string	"Invalid request body or error"
//	@Router			/teams/users [put]
func (tc *TeamController) AddUserToTeam(c *gin.Context) {
	var req dto.UserToTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidRequestBodyError})
		return
	}
	user, team, err := tc.teamService.AddUserToTeam(req.UserID, req.TeamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := dto.NewAddUserToTeamResponse(*user, *team)
	c.JSON(http.StatusOK, resp)
}

// DeleteUserFromTeam
//
//	@Summary		Delete a user from a team
//	@Description	Deletes a user from a team by providing team ID
//
//	@Security		Bearer
//
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserToTeamRequest		true	"User ID and Team ID"
//	@Success		200		{object}	dto.AddUserToTeamResponse	"User removed from team"
//	@Failure		400		{object}	map[string]string			"Invalid request body or error"
//	@Router			/teams/users [delete]
func (tc *TeamController) DeleteUserFromTeam(c *gin.Context) {
	var req dto.UserToTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": InvalidRequestBodyError})
		return
	}
	if req.TeamID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": EmptyParametersError})
		return
	}
	user, team, err := tc.teamService.DeleteUserFromTeam(req.UserID, req.TeamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := dto.NewAddUserToTeamResponse(*user, *team)
	c.JSON(http.StatusOK, resp)
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
//	@Security		Bearer
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
//	@Security		Bearer
//	@Produce		json
//	@Param			id	path		string					true	"Team ID"
//	@Success		200	{object}	map[string]interface{}	"Team deleted"
//	@Failure		400	{object}	map[string]interface{}	"Bad Request: Missing team ID"
//	@Failure		500	{object}	map[string]interface{}	"Internal Server Error"
//	@Router			/teams/{id} [delete]
func (tc *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": MissingTeamIDError})
		return
	}

	if err := tc.teamService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": TeamDeletedMessage})
}
