package service

import (
	"fmt"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

const (
	fileIDEmpty      = "file id can not be empty"
	userNotInTeamErr = "user is not a member of this team"
	teamNotFoundErr  = "team not found"
	fileNotInTeamErr = "file does not belong to this team"
)

type FileServiceInterface interface {
	CreateFile(request *dto.FileUploadRequest, userID string) (*dto.FileUploadResponse, error)
	GetFileByID(id, userID string) (*entity.File, error)
	GetFilesByTeam(teamID, userID string, page, limit int) (*dto.FileListResponse, error)
	DeleteFile(id, userID string) error
}

type FileService struct {
	fileRepo persistence.FileRepositoryInterface
	userRepo UserRepositoryInterface
	teamRepo TeamRepositoryInterface
}

func NewFileService() *FileService {
	return &FileService{
		fileRepo: persistence.NewFileRepository(),
		userRepo: persistence.NewUserRepository(),
		teamRepo: persistence.NewTeamRepository(),
	}
}

func NewFileServiceWithRepo(fileRepo persistence.FileRepositoryInterface, userRepo UserRepositoryInterface, teamRepo TeamRepositoryInterface) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// isUserInTeam checks if user is a member of the specified team
func (fs *FileService) isUserInTeam(userID, teamID string) error {
	user, err := fs.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.TeamsIds == nil {
		return fmt.Errorf(userNotInTeamErr)
	}

	for _, tid := range *user.TeamsIds {
		if tid == teamID {
			return nil
		}
	}
	return fmt.Errorf(userNotInTeamErr)
}

func (fs *FileService) CreateFile(request *dto.FileUploadRequest, userID string) (*dto.FileUploadResponse, error) {
	if err := validator.ValidateFileUpload(request); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify user is member of the team
	if request.ContextType == entity.FileContextTeam {
		if err := fs.isUserInTeam(userID, request.ContextID); err != nil {
			return nil, err
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	file := entity.NewFile(id, request.Name, request.Type, request.Extension, request.Content, request.OwnerID, request.ContextType, request.ContextID, request.Size, now, now)

	if err := fs.fileRepo.Create(file); err != nil {
		return nil, err
	}

	resp := &dto.FileUploadResponse{
		ID:          file.ID,
		Name:        file.Name,
		Type:        file.Type,
		Extension:   file.Extension,
		Size:        file.Size,
		OwnerID:     file.OwnerID,
		ContextType: file.ContextType,
		ContextID:   file.ContextID,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	}
	return resp, nil
}

func (fs *FileService) GetFileByID(id, userID string) (*entity.File, error) {
	if id == "" {
		return nil, fmt.Errorf(fileIDEmpty)
	}

	file, err := fs.fileRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify user has access to this file's context
	if file.ContextType == entity.FileContextTeam {
		if err := fs.isUserInTeam(userID, file.ContextID); err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (fs *FileService) GetFilesByTeam(teamID, userID string, page, limit int) (*dto.FileListResponse, error) {
	// Verify user is member of the team
	if err := fs.isUserInTeam(userID, teamID); err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	files, err := fs.fileRepo.GetByContextID(entity.FileContextTeam, teamID)
	if err != nil {
		return nil, err
	}

	totalCount := len(files)
	totalPages := (totalCount + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}

	paginatedFiles := files[start:end]
	result := make([]*dto.FileUploadResponse, 0, len(paginatedFiles))
	for _, f := range paginatedFiles {
		result = append(result, &dto.FileUploadResponse{
			ID:          f.ID,
			Name:        f.Name,
			Type:        f.Type,
			Extension:   f.Extension,
			Size:        f.Size,
			OwnerID:     f.OwnerID,
			ContextType: f.ContextType,
			ContextID:   f.ContextID,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		})
	}

	return &dto.FileListResponse{
		Files:      result,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (fs *FileService) DeleteFile(id, userID string) error {
	if id == "" {
		return fmt.Errorf(fileIDEmpty)
	}

	file, err := fs.fileRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verify user has access (is in the team)
	if file.ContextType == entity.FileContextTeam {
		if err := fs.isUserInTeam(userID, file.ContextID); err != nil {
			return err
		}
	}

	return fs.fileRepo.Delete(id)
}
