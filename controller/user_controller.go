package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

const (
	userNotFoundError           = "User not found"
	userDeletedSuccessfully     = "User deleted successfully"
	statisticsUpdatedSuccessfully = "Statistics updated successfully"
)

type UserController struct {
	userService UserServiceInterface
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func NewUserControllerWithService(userService UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

type UserServiceInterface interface {
	SignUp(request *dto.SignUpUserRequest) (*dto.SignUpUserResponse, error)
	GetUserByID(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id string) error
	GetAllUsers() ([]*entity.User, error)
	UpdateUserStatistics(id string, timeSpentOnApp int64, timeSpentOnTeam model.TimeSpentOnTeam) (*entity.User, error)
}

// SignUp
//
//	@Summary	Register a new user
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.SignUpUserRequest	true	"The sign-up request"
//	@Success	201		{object}	dto.SignUpUserResponse
//	@Router		/users/signup [post]
func (uc *UserController) SignUp(c *gin.Context) {
	var request dto.SignUpUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userService.SignUp(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser
//
//	@Summary	Get a user by ID
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The user's ID"
//	@Success	200	{object}	entity.User
//	@Router		/users/{id}  [get]
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": userNotFoundError})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers
//
//	@Summary	Get all users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]entity.User
//	@Router		/users [get]
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser TODO: Add Swagger comment
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": userNotFoundError})
		return
	}

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser
//
//	@Summary	Delete a user
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"The user's ID"
//	@Success	200
//	@Router		/users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := uc.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userDeletedSuccessfully})
}

// UpdateUserStatistics
//
//	@Summary	Update user statistics
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"The user's ID"
//	@Param		request	body		dto.UpdateStatisticsRequest	true	"The statistics update request"
//	@Success	200		{object}	dto.UpdateStatisticsResponse
//	@Router		/users/{id}/statistics [put]
func (uc *UserController) UpdateUserStatistics(c *gin.Context) {
	id := c.Param("id")

	var request dto.UpdateStatisticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamTimeSpent := model.TimeSpentOnTeam{
		TeamId:   request.TeamId,
		Duration: request.TimeSpentOnTeam,
	}

	updatedUser, err := uc.userService.UpdateUserStatistics(id, request.TimeSpentOnApp, teamTimeSpent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.NewUpdateStatisticsResponse(updatedUser.ID, updatedUser.Statistics)
	c.JSON(http.StatusOK, response)
}
