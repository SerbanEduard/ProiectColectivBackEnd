package tests

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]*entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUp(request *dto.SignUpUserRequest) (*dto.SignUpUserResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SignUpUserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers() ([]*entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserService) UpdateUserStatistics(id string, timeSpentOnApp int64, timeSpentOnTeam model.TimeSpentOnTeam) (*entity.User, error) {
	args := m.Called(id, timeSpentOnApp, timeSpentOnTeam)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockFriendRequestRepository struct {
	mock.Mock
}

func (m *MockFriendRequestRepository) Create(request *entity.FriendRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) GetByUsers(fromUserID, toUserID string) (*entity.FriendRequest, error) {
	args := m.Called(fromUserID, toUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FriendRequest), args.Error(1)
}

func (m *MockFriendRequestRepository) Update(request *entity.FriendRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) GetPendingRequestsForUser(userID string) ([]*entity.FriendRequest, error) {
	args := m.Called(userID)
	return args.Get(0).([]*entity.FriendRequest), args.Error(1)
}

type MockFriendRequestService struct {
	mock.Mock
}

func (m *MockFriendRequestService) SendFriendRequest(fromUserID, toUserID string) error {
	args := m.Called(fromUserID, toUserID)
	return args.Error(0)
}

func (m *MockFriendRequestService) RespondToFriendRequest(fromUserID, toUserID string, accept bool) error {
	args := m.Called(fromUserID, toUserID, accept)
	return args.Error(0)
}

func (m *MockFriendRequestService) GetPendingRequests(userID string) ([]*entity.FriendRequest, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FriendRequest), args.Error(1)
}

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(team *entity.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetTeamById(id string) (*entity.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error) {
	args := m.Called(prefix, x)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeamsByName(name string) ([]*entity.Team, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) GetAll() ([]*entity.Team, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Team), args.Error(1)
}

func (m *MockTeamRepository) Update(team *entity.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}