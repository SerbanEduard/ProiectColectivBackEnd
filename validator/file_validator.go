package validator

import (
	"errors"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

// MaxFileSize is the maximum allowed file size in bytes (500 MB)
const MaxFileSize int64 = 500 * 1024 * 1024

func ValidateFileUpload(req *dto.FileUploadRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Type) == "" {
		return errors.New("type is required")
	}
	if strings.TrimSpace(req.Extension) == "" {
		return errors.New("extension is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content is required")
	}
	if strings.TrimSpace(req.OwnerID) == "" {
		return errors.New("ownerId is required")
	}
	if req.Size <= 0 {
		return errors.New("size must be greater than zero")
	}
	if req.Size > MaxFileSize {
		return errors.New("file size exceeds maximum allowed (500 MB)")
	}

	// Validate context
	if strings.TrimSpace(req.ContextType) == "" {
		return errors.New("contextType is required")
	}
	if req.ContextType != entity.FileContextTeam && req.ContextType != entity.FileContextChat {
		return errors.New("contextType must be 'team' or 'chat'")
	}
	if strings.TrimSpace(req.ContextID) == "" {
		return errors.New("contextId is required")
	}

	if after, ok := strings.CutPrefix(req.Extension, "."); ok {
		req.Extension = after
	}

	return nil
}
