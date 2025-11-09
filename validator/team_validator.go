package validator

import "github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"

func ValidateTeamRequest(request *dto.TeamRequest) error {
	validations := []func() error{
		func() error { return validateRequired(request.Name, "first name is required") },
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
