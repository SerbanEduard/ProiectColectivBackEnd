package service_test

import (
	"fmt"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	. "github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	TestDuration2Hour  = int64(7200000)  // 2 hours in milliseconds
	TestDuration1Hour  = int64(3600000)  // 1 hour in milliseconds
	TestDuration30Min  = int64(1800000)  // 30 minutes in milliseconds
	TestDuration75Min  = int64(4500000)  // 75 minutes in milliseconds
	TestDuration45Min  = int64(2700000)  // 45 minutes in milliseconds
	TestDuration3Hour  = int64(10800000) // 3 hours in milliseconds
	TestDuration105Min = int64(6300000)  // 105 minutes in milliseconds
)

func TestUserService_SignUp_Success(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	request := &tests.ValidSignUpRequest

	mockRepo.On("GetByUsername", TestUsername).Return(nil, fmt.Errorf(ErrUserNotFound))
	mockRepo.On("GetByEmail", TestEmail).Return(nil, fmt.Errorf(ErrUserNotFound))
	mockRepo.On("Create", mock.MatchedBy(func(user *entity.User) bool {
		return user.FirstName == TestFirstName &&
			user.LastName == TestLastName &&
			user.Username == TestUsername &&
			user.Email == TestEmail &&
			user.TopicsOfInterest != nil &&
			len(*user.TopicsOfInterest) == 1 &&
			(*user.TopicsOfInterest)[0] == model.Programming &&
			user.ID != "" &&
			user.Password != TestPassword
	})).Return(nil)

	response, err := userService.SignUp(request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, TestFirstName, response.FirstName)
	assert.Equal(t, TestLastName, response.LastName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_SignUp_UsernameExists(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	request := &tests.ExistingUsernameRequest

	existingUser := &tests.ExistingUser
	mockRepo.On("GetByUsername", ExistingUsername).Return(existingUser, nil)

	response, err := userService.SignUp(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrUsernameExists, err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_SignUp_EmailExists(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	request := &tests.ExistingEmailRequest

	mockRepo.On("GetByUsername", TestUsername).Return(nil, fmt.Errorf(ErrUserNotFound))
	mockRepo.On("GetByEmail", ExistingEmail).Return(&tests.ExistingUser, nil)

	response, err := userService.SignUp(request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrEmailExists, err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserStatistics_Success(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	user := &entity.User{
		ID: TestUserID,
		Statistics: &model.Statistics{
			TotalTimeSpentOnApp: TestDuration1Hour,
			TimeSpentOnTeams: []model.TimeSpentOnTeam{
				{TeamId: TestTeamID, Duration: TestDuration30Min},
			},
		},
	}

	mockRepo.On("GetByID", TestUserID).Return(user, nil)
	mockTeamRepo.On("GetTeamById", TestTeamID).Return(&entity.Team{Id: TestTeamID}, nil)

	expectedStatistics := dto.NewStatisticsResponse(TestUserID, user.Statistics)

	statistics, err := userService.GetUserStatistics(TestUserID)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatistics, statistics)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserStatistics_Success(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	user := &entity.User{
		ID: TestUserID,
		Statistics: &model.Statistics{
			TotalTimeSpentOnApp: TestDuration1Hour,
			TimeSpentOnTeams: []model.TimeSpentOnTeam{
				{TeamId: TestTeamID, Duration: TestDuration30Min},
			},
		},
	}

	mockRepo.On("GetByID", TestUserID).Return(user, nil)
	mockTeamRepo.On("GetTeamById", TestTeamID).Return(&entity.Team{Id: TestTeamID}, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *entity.User) bool {
		return u.Statistics.TotalTimeSpentOnApp == TestDuration3Hour &&
			len(u.Statistics.TimeSpentOnTeams) == 1 &&
			u.Statistics.TimeSpentOnTeams[0].Duration == TestDuration105Min
	})).Return(nil)

	timeSpentOnTeam := model.TimeSpentOnTeam{
		TeamId:   TestTeamID,
		Duration: TestDuration75Min,
	}

	_, err := userService.UpdateUserStatistics(TestUserID, TestDuration2Hour, timeSpentOnTeam)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserStatistics_NewTeam(t *testing.T) {
	mockRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	userService := service.NewUserServiceWithRepo(mockRepo, mockTeamRepo)

	user := &entity.User{
		ID:         TestUserID,
		Statistics: nil,
	}

	mockRepo.On("GetByID", TestUserID).Return(user, nil)
	mockTeamRepo.On("GetTeamById", TestTeamID2).Return(&entity.Team{Id: TestTeamID2}, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *entity.User) bool {
		return u.Statistics != nil &&
			u.Statistics.TotalTimeSpentOnApp == TestDuration1Hour &&
			len(u.Statistics.TimeSpentOnTeams) == 1 &&
			u.Statistics.TimeSpentOnTeams[0].TeamId == TestTeamID2
	})).Return(nil)

	timeSpentOnTeam := model.TimeSpentOnTeam{
		TeamId:   TestTeamID2,
		Duration: TestDuration45Min,
	}

	_, err := userService.UpdateUserStatistics(TestUserID, TestDuration1Hour, timeSpentOnTeam)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
