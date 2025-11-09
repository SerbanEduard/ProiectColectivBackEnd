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

	ValidFriendRequest = entity.FriendRequest{
        FromUserID: "user1",
        ToUserID:   "user2",
        Status:     entity.PENDING,
    }

    AcceptedFriendRequest = entity.FriendRequest{
        FromUserID: "user1",
        ToUserID:   "user2",
        Status:     entity.ACCEPTED,
    }

    DeniedFriendRequest = entity.FriendRequest{
        FromUserID: "user1",
        ToUserID:   "user2",
        Status:     entity.DENIED,
    }
)