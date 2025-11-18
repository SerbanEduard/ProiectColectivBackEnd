package dto

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

type SignUpUserRequest struct {
	FirstName        string                   `json:"firstname"`
	LastName         string                   `json:"lastname"`
	Username         string                   `json:"username"`
	Email            string                   `json:"email"`
	Password         string                   `json:"password"`
	TopicsOfInterest *[]model.TopicOfInterest `json:"topicsOfInterest,omitempty"`
}

type SignUpUserResponse struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
}

func NewSignUpUserRequest(firstName, lastName, username, email, password string, topicsOfInterest *[]model.TopicOfInterest) *SignUpUserRequest {
	return &SignUpUserRequest{
		FirstName:        firstName,
		LastName:         lastName,
		Username:         username,
		Email:            email,
		Password:         password,
		TopicsOfInterest: topicsOfInterest,
	}
}

func NewSignUpUserResponse(firstName, lastName, username string) *SignUpUserResponse {
	return &SignUpUserResponse{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}
}

type AddStatisticsToUserRequest struct {
	Username   string           `json:"username"`
	Statistics model.Statistics `json:"statistics"`
}

func NewAddStatisticsToUserRequest(username string, statistics model.Statistics) *AddStatisticsToUserRequest {
	return &AddStatisticsToUserRequest{
		Username:   username,
		Statistics: statistics,
	}
}

// Login DTOs
type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	TokenType   string       `json:"tokenType"`
	ExpiresIn   string       `json:"expiresIn"`
	User        UserResponse `json:"user"`
}

// UserResponse is a safe representation of the user returned to clients (no password).
type UserResponse struct {
	ID               string                   `json:"id"`
	FirstName        string                   `json:"firstname"`
	LastName         string                   `json:"lastname"`
	Username         string                   `json:"username"`
	Email            string                   `json:"email"`
	TopicsOfInterest *[]model.TopicOfInterest `json:"topicsOfInterest,omitempty"`
	TeamsIds         *[]string                `json:"teams,omitempty"`
	Statistics       *model.Statistics        `json:"statistics,omitempty"`
}

// NewUserResponse converts an entity.User to a safe UserResponse (omits password).
func NewUserResponse(u *entity.User) UserResponse {
	resp := UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
	}
	if u.TopicsOfInterest != nil {
		resp.TopicsOfInterest = u.TopicsOfInterest
	} else {
		empty := []model.TopicOfInterest{}
		resp.TopicsOfInterest = &empty
	}
	if u.TeamsIds != nil {
		resp.TeamsIds = u.TeamsIds
	}
	if u.Statistics != nil {
		resp.Statistics = u.Statistics
	}
	return resp
}

func NewLoginResponse(token, expiresIn string, user *entity.User) *LoginResponse {
	resp := &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User:        NewUserResponse(user),
	}
	return resp
}

// UserUpdateRequestDTO is used for updating user profile (all fields optional)
type UserUpdateRequestDTO struct {
	FirstName        string                   `json:"firstname,omitempty"`
	LastName         string                   `json:"lastname,omitempty"`
	Username         string                   `json:"username,omitempty"`
	Email            string                   `json:"email,omitempty"`
	TopicsOfInterest *[]model.TopicOfInterest `json:"topicsOfInterest,omitempty"`
}

// UserUpdateResponseDTO is a safe representation after update (no password)
type UserUpdateResponseDTO struct {
	ID               string                   `json:"id"`
	FirstName        string                   `json:"firstname"`
	LastName         string                   `json:"lastname"`
	Username         string                   `json:"username"`
	Email            string                   `json:"email"`
	TopicsOfInterest *[]model.TopicOfInterest `json:"topicsOfInterest,omitempty"`
	TeamsIds         *[]string                `json:"teams,omitempty"`
	Statistics       *model.Statistics        `json:"statistics,omitempty"`
}

// NewUserUpdateResponseDTO converts an entity.User to UserUpdateResponseDTO
func NewUserUpdateResponseDTO(u *entity.User) *UserUpdateResponseDTO {
	resp := &UserUpdateResponseDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
	}
	if u.TopicsOfInterest != nil {
		resp.TopicsOfInterest = u.TopicsOfInterest
	} else {
		empty := []model.TopicOfInterest{}
		resp.TopicsOfInterest = &empty
	}
	if u.TeamsIds != nil {
		resp.TeamsIds = u.TeamsIds
	}
	if u.Statistics != nil {
		resp.Statistics = u.Statistics
	}
	return resp
}

// UserPasswordRequestDTO is used for password updates (requires old password verification)
type UserPasswordRequestDTO struct {
	ID          string `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
