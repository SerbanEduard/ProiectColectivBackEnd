package service_test

import (
	"testing"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileService_CreateFile_Success(t *testing.T) {
	mockFileRepo := new(tests.MockFileRepository)
	mockUserRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	fs := service.NewFileServiceWithRepo(mockFileRepo, mockUserRepo, mockTeamRepo)

	userID := "user1"
	teamID := "team1"
	teamIDs := []string{teamID}

	req := &dto.FileUploadRequest{
		Name:        "test.txt",
		Type:        "text/plain",
		Extension:   "txt",
		Content:     "dGVzdA==",
		OwnerID:     userID,
		Size:        4,
		ContextType: entity.FileContextTeam,
		ContextID:   teamID,
	}

	mockUserRepo.On("GetByID", userID).Return(&entity.User{ID: userID, TeamsIds: &teamIDs}, nil)
	mockFileRepo.On("Create", mock.MatchedBy(func(f *entity.File) bool {
		return f.Name == req.Name && f.OwnerID == req.OwnerID && f.ContextType == entity.FileContextTeam && f.ContextID == teamID
	})).Return(nil)

	resp, err := fs.CreateFile(req, userID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, entity.FileContextTeam, resp.ContextType)
	assert.Equal(t, teamID, resp.ContextID)

	mockFileRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestFileService_CreateFile_UserNotInTeam(t *testing.T) {
	mockFileRepo := new(tests.MockFileRepository)
	mockUserRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	fs := service.NewFileServiceWithRepo(mockFileRepo, mockUserRepo, mockTeamRepo)

	userID := "user1"
	otherTeamIDs := []string{"other-team"}

	req := &dto.FileUploadRequest{
		Name:        "test.txt",
		Type:        "text/plain",
		Extension:   "txt",
		Content:     "dGVzdA==",
		OwnerID:     userID,
		Size:        4,
		ContextType: entity.FileContextTeam,
		ContextID:   "team1",
	}

	mockUserRepo.On("GetByID", userID).Return(&entity.User{ID: userID, TeamsIds: &otherTeamIDs}, nil)

	resp, err := fs.CreateFile(req, userID)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "not a member")

	mockFileRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestFileService_GetFilesByTeam_Success(t *testing.T) {
	mockFileRepo := new(tests.MockFileRepository)
	mockUserRepo := new(tests.MockUserRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	fs := service.NewFileServiceWithRepo(mockFileRepo, mockUserRepo, mockTeamRepo)

	userID := "user1"
	teamID := "team1"
	teamIDs := []string{teamID}
	now := time.Now().Unix()

	files := []*entity.File{{
		ID: "1", Name: "a.txt", Type: "text/plain", Extension: "txt",
		Size: 1, OwnerID: userID, ContextType: entity.FileContextTeam, ContextID: teamID,
		CreatedAt: now, UpdatedAt: now,
	}}

	mockUserRepo.On("GetByID", userID).Return(&entity.User{ID: userID, TeamsIds: &teamIDs}, nil)
	mockFileRepo.On("GetByContextID", entity.FileContextTeam, teamID).Return(files, nil)

	got, err := fs.GetFilesByTeam(teamID, userID, 1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got.Files, 1)
	assert.Equal(t, "1", got.Files[0].ID)
	assert.Equal(t, entity.FileContextTeam, got.Files[0].ContextType)

	mockFileRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
