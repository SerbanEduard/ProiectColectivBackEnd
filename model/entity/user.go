package entity

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type User struct {
	ID               string                  `json:"id"`
	FirstName        string                  `json:"firstname"`
	LastName         string                  `json:"lastname"`
	Username         string                  `json:"username"`
	Email            string                  `json:"email"`
	Password         string                  `json:"password"`
	TopicsOfInterest []model.TopicOfInterest `json:"topicsOfInterest"`
	TeamsIds         []string                `json:"teams"`
}

func NewUser(id, firstName, lastName, username, email, password string, topicsOfInterest []model.TopicOfInterest) *User {
	return &User{
		ID:               id,
		FirstName:        firstName,
		LastName:         lastName,
		Username:         username,
		Email:            email,
		Password:         password,
		TopicsOfInterest: topicsOfInterest,
	}
}
