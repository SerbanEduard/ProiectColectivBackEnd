package controller

import (
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	userNotFoundError             = "User not found"
	userDeletedSuccessfully       = "User deleted successfully"
	statisticsUpdatedSuccessfully = "Statistics updated successfully"
	invalidTimeSpentOnAppFormat   = "Invalid timeSpentOnApp format"
	invalidTimeSpentOnTeamFormat  = "Invalid timeSpentOnTeam format"
	invalidCredentials            = "invalid email or password"
	jwtExpiresHours               = 24
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

// Login
//
//	@Summary    Login user and return JWT
//	@Accept     json
//	@Produce    json
//	@Param      request body        dto.LoginRequest true "The login request"
//	@Success    200     {object}    dto.LoginResponse
//	@Router     /users/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user *entity.User
	var err error
	// allow login by email or username
	if req.Email != "" {
		user, err = uc.userService.GetUserByEmail(req.Email)
	} else if req.Username != "" {
		user, err = uc.userService.GetUserByUsername(req.Username)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username required"})
		return
	}
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidCredentials})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidCredentials})
		return
	}

	secret := config.GetJWTSecret()
	expiration := time.Now().Add(time.Hour * jwtExpiresHours)

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiration.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.NewLoginResponse(signed, "24h", user)
	c.JSON(http.StatusOK, resp)
}
