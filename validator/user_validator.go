package validator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
)

func ValidateSignUpRequest(request *dto.SignUpUserRequest) error {
	validations := []func() error{
		func() error { return validateRequired(request.FirstName, "first name is required") },
		func() error { return validateRequired(request.LastName, "last name is required") },
		func() error { return validateRequired(request.Username, "username is required") },
		func() error { return validateMinLength(request.Username, 3, "username must be at least 3 characters") },
		func() error { return validateRequired(request.Email, "email is required") },
		func() error { return validateEmail(request.Email) },
		func() error { return validateMinLength(request.Password, 6, "password must be at least 6 characters") },
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}

func validateRequired(value, message string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(message)
	}
	return nil
}

func validateMinLength(value string, minLen int, message string) error {
	if len(value) < minLen {
		return errors.New(message)
	}
	return nil
}

func validateEmail(email string) error {
	if !IsValidEmail(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
