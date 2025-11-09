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

func (tc *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")
	team, err := tc.teamService.GetTeamById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (tc *TeamController) GetAllTeams(c *gin.Context) {
	teams, err := tc.teamService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

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
