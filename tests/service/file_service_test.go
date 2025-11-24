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
	mockRepo := new(tests.MockFileRepository)
	fs := service.NewFileServiceWithRepo(mockRepo)

	req := &dto.FileUploadRequest{
		Name:      "test.txt",
		Type:      "text/plain",
		Extension: "txt",
		Content:   "dGVzdA==", // "test" base64
		OwnerID:   "owner1",
		Size:      4,
	}

	// Expect Create called with any *entity.File and return nil
	mockRepo.On("Create", mock.MatchedBy(func(f *entity.File) bool {
		return f.Name == req.Name && f.OwnerID == req.OwnerID && f.Size == req.Size
	})).Return(nil)

	resp, err := fs.CreateFile(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Name)

	mockRepo.AssertExpectations(t)
}

func TestFileService_CreateFile_ValidationFail(t *testing.T) {
	mockRepo := new(tests.MockFileRepository)
	fs := service.NewFileServiceWithRepo(mockRepo)

	req := &dto.FileUploadRequest{
		Name: "",
	}

	resp, err := fs.CreateFile(req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Ensure repo not called
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestFileService_GetAllFiles(t *testing.T) {
	mockRepo := new(tests.MockFileRepository)
	fs := service.NewFileServiceWithRepo(mockRepo)

	now := time.Now().Unix()
	files := []*entity.File{{ID: "1", Name: "a.txt", Size: 1, OwnerID: "o", CreatedAt: now}}
	mockRepo.On("GetAll").Return(files, nil)

	got, err := fs.GetAllFiles()
	assert.NoError(t, err)
	assert.Equal(t, files, got)

	mockRepo.AssertExpectations(t)
}
