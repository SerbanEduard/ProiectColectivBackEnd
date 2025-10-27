package tests

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

var (
	ValidSignUpRequest = dto.SignUpUserRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Username:         "johndoe",
		Email:            "john@example.com",
		Password:         "password123",
		TopicsOfInterest: []model.TopicOfInterest{model.Programming},
	}

	ExistingUsernameRequest = dto.SignUpUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "existinguser",
		Email:     "john@example.com",
		Password:  "password123",
	}

	ExistingEmailRequest = dto.SignUpUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "existing@example.com",
		Password:  "password123",
	}

	ValidSignUpResponse = dto.SignUpUserResponse{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
	}

	ExistingUser = entity.User{
		Username: "existinguser",
	}

	// Login test data
	ValidLoginRequest = dto.LoginUserRequest{
		Email:    "john@example.com",
		Password: "password123",
	}

	ValidLoginResponse = dto.LoginUserResponse{
		ID:        "user123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
	}

	ValidUserWithPassword = entity.User{
		ID:        "user123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Password:  "$2a$10$abcdefghijklmnopqrstuvwxyzABCDEF", // bcrypt hash pentru "password123"
	}
)