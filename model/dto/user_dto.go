package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   string `json:"expiresIn"`
	User        struct {
		Id               string                   `json:"id"`
		Username         string                   `json:"username"`
		Email            string                   `json:"email"`
		TopicsOfInterest *[]model.TopicOfInterest `json:"topicsOfInterest,omitempty"`
	} `json:"user"`
}

func NewLoginResponse(token, expiresIn, userId, username, email string, topics *[]model.TopicOfInterest) *LoginResponse {
	resp := &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}
	resp.User.Id = userId
	resp.User.Username = username
	resp.User.Email = email
	// ensure we return an empty slice instead of omitting the field when topics is nil
	if topics == nil {
		empty := []model.TopicOfInterest{}
		resp.User.TopicsOfInterest = &empty
	} else {
		resp.User.TopicsOfInterest = topics
	}
	return resp
}
