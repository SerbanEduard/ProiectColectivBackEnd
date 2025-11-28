package validator

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
)

func ValidateDirectMessageRequest(request *dto.DirectMessageRequest) error {
	validations := []func() error{
		func() error { return validateRequired(request.SenderID, "sender id is required") },
		func() error { return validateRequired(request.ReceiverID, "receiver(user) id is required") },
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}

func ValidateTeamMessageRequest(request *dto.TeamMessageRequest) error {
	validations := []func() error{
		func() error { return validateRequired(request.SenderID, "sender id is required") },
		func() error { return validateRequired(request.TeamId, "team id is required") },
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
