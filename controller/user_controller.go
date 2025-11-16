package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

const (
	userNotFoundError       = "User not found"
	userDeletedSuccessfully = "User deleted successfully"
	invalidCredentials      = "invalid email or password"
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
	GetUserByUsername(username string) (*entity.User, error)
	Login(request *dto.LoginRequest) (*dto.LoginResponse, error)
	UpdateUser(user *entity.User) error
	UpdateUserProfile(userID string, req *dto.UserUpdateRequestDTO) (*dto.UserUpdateResponseDTO, error)
	UpdateUserPassword(userID string, req *dto.UserPasswordRequestDTO) error
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
//	@Failure	400		{object}	map[string]string
//	@Failure	409		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/users/signup [post]
func (uc *UserController) SignUp(c *gin.Context) {
	var request dto.SignUpUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userService.SignUp(&request)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") ||
			strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "required") ||
			strings.Contains(err.Error(), "must") {
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		// All other errors are server errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser
//
//	@Summary	Get a user by ID
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The user's ID"
//	@Success	200	{object}	entity.User
//	@Failure	404	{object}	map[string]string
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
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]entity.User
//	@Failure	500	{object}	map[string]string
//	@Router		/users [get]
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser
//
//	@Summary	Update user profile (selective fields)
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"The user's ID"
//	@Param		request	body		dto.UserUpdateRequestDTO	true	"The user profile update (all fields optional)"
//	@Success	200		{object}	dto.UserUpdateResponseDTO
//	@Failure	400		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/users/{id} [patch]
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req dto.UserUpdateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.userService.UpdateUserProfile(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserPassword
//
//	@Summary	Update user password
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"The user's ID"
//	@Param		request	body		dto.UserPasswordRequestDTO	true	"The password update request"
//	@Success	200		{object}	map[string]string
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/users/{id}/password [put]
func (uc *UserController) UpdateUserPassword(c *gin.Context) {
	id := c.Param("id")

	var req dto.UserPasswordRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	if err := uc.userService.UpdateUserPassword(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// DeleteUser
//
//	@Summary	Delete a user
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The user's ID"
//	@Success	200	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := uc.userService.DeleteUser(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userDeletedSuccessfully})
}

// UpdateUserStatistics
//
//	@Summary	Update user statistics
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"The user's ID"
//	@Param		request	body		dto.UpdateStatisticsRequest	true	"The statistics update request"
//	@Success	200		{object}	dto.UpdateStatisticsResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
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
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.NewUpdateStatisticsResponse(updatedUser.ID, updatedUser.Statistics)
	c.JSON(http.StatusOK, response)
}

// Login
//
//	@Summary		Login user by email or username and return JWT
//	@Description	Accepts either `email` or `username` along with `password`. Returns an access token and the full user (without password).
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"The login request (email or username + password)"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.userService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": invalidCredentials})
			return
		}
		if strings.Contains(err.Error(), "required") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
