package validator

import (
	"errors"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
)

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

	if after, ok := strings.CutPrefix(req.Extension, "."); ok {
		req.Extension = after
	}

	return nil
}
