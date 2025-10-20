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
