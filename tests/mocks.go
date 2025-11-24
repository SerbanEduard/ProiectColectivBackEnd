package tests

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/stretchr/testify/mock"
)

type QuizRepositoryInterface interface {
	Create(quiz entity.Quiz) error
	Update(quiz entity.Quiz) error
	GetById(id string) (entity.Quiz, error)
	GetByUser(id string) ([]entity.Quiz, error)
	GetByTeam(id string) ([]entity.Quiz, error)
}
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

func (m *MockUserService) GetUserByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) Login(request *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
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

func (m *MockUserService) GetUserStatistics(id string) (*dto.StatisticsResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StatisticsResponse), args.Error(1)
}

func (m *MockUserService) UpdateUserStatistics(id string, timeSpentOnApp int64, timeSpentOnTeam model.TimeSpentOnTeam) (*entity.User, error) {
	args := m.Called(id, timeSpentOnApp, timeSpentOnTeam)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) UpdateUserProfile(userID string, req *dto.UserUpdateRequestDTO) (*dto.UserUpdateResponseDTO, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserUpdateResponseDTO), args.Error(1)
}

func (m *MockUserService) UpdateUserPassword(userID string, req *dto.UserPasswordRequestDTO) error {
	args := m.Called(userID, req)
	return args.Error(0)
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

func (m *MockFriendRequestRepository) GetFriendsForUser(userID string) ([]string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
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

func (m *MockFriendRequestService) GetFriends(userID string) ([]*entity.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockFriendRequestService) GetMutualFriends(userA, userB string) ([]*entity.User, error) {
	args := m.Called(userA, userB)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
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

type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) Create(quiz entity.Quiz) error {
	args := m.Called(quiz)
	return args.Error(0)
}

func (m *MockQuizRepository) Update(quiz entity.Quiz) error {
	args := m.Called(quiz)
	return args.Error(0)
}

func (m *MockQuizRepository) GetById(id string) (entity.Quiz, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		// Asigură-te că primul argument este o valoare zero (entity.Quiz{}) dacă este nil, și returnează eroarea.
		return entity.Quiz{}, args.Error(1)
	}
	return args.Get(0).(entity.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetByUser(id string) ([]entity.Quiz, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetByTeam(id string) ([]entity.Quiz, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Quiz), args.Error(1)
}

// MockFileRepository is used for file service tests
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(file *entity.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(id string) (*entity.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.File), args.Error(1)
}

func (m *MockFileRepository) GetAll() ([]*entity.File, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.File), args.Error(1)
}

func (m *MockFileRepository) Update(file *entity.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockQuizService is used by controller tests to mock service layer behavior.
type MockQuizService struct {
	mock.Mock
}

func (m *MockQuizService) CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error) {
	args := m.Called(request)
	var resp dto.CreateQuizResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(dto.CreateQuizResponse)
	}
	return resp, args.Error(1)
}

func (m *MockQuizService) GetQuizWithAnswersById(id string) (entity.Quiz, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return entity.Quiz{}, args.Error(1)
	}
	return args.Get(0).(entity.Quiz), args.Error(1)
}

func (m *MockQuizService) GetQuizWithoutAnswersById(id string) (dto.ReadQuizResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return dto.ReadQuizResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.ReadQuizResponse), args.Error(1)
}

func (m *MockQuizService) SolveQuiz(request dto.SolveQuizRequest, userId string) (dto.SolveQuizResponse, error) {
	args := m.Called(request, userId)
	if args.Get(0) == nil {
		return dto.SolveQuizResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.SolveQuizResponse), args.Error(1)
}
