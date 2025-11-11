package tests

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

var (
	ValidSignUpRequest = dto.SignUpUserRequest{
		FirstName:        TestFirstName,
		LastName:         TestLastName,
		Username:         TestUsername,
		Email:            TestEmail,
		Password:         TestPassword,
		TopicsOfInterest: &[]model.TopicOfInterest{model.Programming},
	}

	ExistingUsernameRequest = dto.SignUpUserRequest{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  ExistingUsername,
		Email:     TestEmail,
		Password:  TestPassword,
	}

	ExistingEmailRequest = dto.SignUpUserRequest{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  TestUsername,
		Email:     ExistingEmail,
		Password:  TestPassword,
	}

	ValidSignUpResponse = dto.SignUpUserResponse{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  TestUsername,
	}

	ExistingUser = entity.User{
		Username: ExistingUsername,
	}

	ValidUpdateStatisticsRequest = dto.UpdateStatisticsRequest{
		TimeSpentOnApp:  TestDurationApp,
		TeamId:          TestTeamID,
		TimeSpentOnTeam: TestDurationTeam,
	}

	ValidTimeSpentOnTeam = model.TimeSpentOnTeam{
		TeamId:   TestTeamID,
		Duration: int64(4500000), // 75 minutes in milliseconds
	}

	ValidFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.PENDING,
	}

	AcceptedFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.ACCEPTED,
	}

	DeniedFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.DENIED,
	}
)
