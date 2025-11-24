package service

import (
	"fmt"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/utils/validator"
)

const (
	fileIDEmpty = "file id can not be empty"
)

type FileServiceInterface interface {
	CreateFile(request *dto.FileUploadRequest) (*dto.FileResponse, error)
	GetFileByID(id string) (*entity.File, error)
	GetAllFiles() ([]*entity.File, error)
	DeleteFile(id string) error
}

type FileService struct {
	fileRepo persistence.FileRepositoryInterface
}

func NewFileService() *FileService {
	return &FileService{
		fileRepo: persistence.NewFileRepository(),
	}
}

func NewFileServiceWithRepo(repo persistence.FileRepositoryInterface) *FileService {
	return &FileService{fileRepo: repo}
}

func (fs *FileService) CreateFile(request *dto.FileUploadRequest) (*dto.FileResponse, error) {
	if err := validator.ValidateFileUpload(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	file := entity.NewFile(id, request.Name, request.Type, request.Extension, request.Content, request.OwnerID, request.Size, now, now)

	if err := fs.fileRepo.Create(file); err != nil {
		return nil, err
	}

	resp := &dto.FileResponse{
		ID:        file.ID,
		Name:      file.Name,
		Type:      file.Type,
		Extension: file.Extension,
		Size:      file.Size,
		OwnerID:   file.OwnerID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
	return resp, nil
}

func (fs *FileService) GetFileByID(id string) (*entity.File, error) {
	if id == "" {
		return nil, fmt.Errorf(fileIDEmpty)
	}
	return fs.fileRepo.GetByID(id)
}

func (fs *FileService) GetAllFiles() ([]*entity.File, error) {
	return fs.fileRepo.GetAll()
}

func (fs *FileService) DeleteFile(id string) error {
	if id == "" {
		return fmt.Errorf(fileIDEmpty)
	}
	return fs.fileRepo.Delete(id)
}
