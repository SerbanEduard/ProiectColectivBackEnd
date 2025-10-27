package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type SignUpUserRequest struct {
	FirstName        string                  `json:"firstname"`
	LastName         string                  `json:"lastname"`
	Username         string                  `json:"username"`
	Email            string                  `json:"email"`
	Password         string                  `json:"password"`
	TopicsOfInterest []model.TopicOfInterest `json:"topicsOfInterest"`
}

type SignUpUserResponse struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
}

func NewSignUpUserRequest(firstName, lastName, username, email, password string, topicsOfInterest []model.TopicOfInterest) *SignUpUserRequest {
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

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
}

func NewLoginUserRequest(email, password string) *LoginUserRequest {
	return &LoginUserRequest{
		Email:    email,
		Password: password,
	}
}

func NewLoginUserResponse(id, firstName, lastName, username string) *LoginUserResponse {
	return &LoginUserResponse{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}
}
