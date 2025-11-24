package validator

import (
	"errors"
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

var (
	ErrValidation = errors.New("validation failed")
)

const (
	nameEmptyError        = "name can not be null"
	invalidQuestionsError = "questions are invalid"
	quizIdEmpty           = "no id specified"
	userIdEmptyError      = "user id cannot be empty"
	teamIdEmptyError      = "team id cannot be empty"
	pageSizeInvalidError  = "page size must be positive"
)

// ValidateCreateQuizRequest validates the quiz creation request
func ValidateCreateQuizRequest(request entity.Quiz) error {
	if request.QuizName == "" {
		return fmt.Errorf("%w: %s", ErrValidation, nameEmptyError)
	}

	for _, question := range request.Questions {
		if question.Question == "" || len(question.Options) == 0 || len(question.Answers) == 0 || question.Type == "" {
			return fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
		}
	}

	return nil
}

// ValidateQuizId validates that quiz ID is not empty
func ValidateQuizId(id string) error {
	if id == "" {
		return fmt.Errorf("%w: %s", ErrValidation, quizIdEmpty)
	}
	return nil
}

// ValidateSolveQuizRequest validates the solve quiz request
func ValidateSolveQuizRequest(request dto.SolveQuizRequest, questions []entity.Question, quizId string) error {
	if quizId == "" {
		return fmt.Errorf("%w: %s", ErrValidation, quizIdEmpty)
	}

	questionsSubmitted := request.Attempts
	if len(questionsSubmitted) != len(questions) {
		return fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
	}

	return nil
}

// ValidateQuestionSubmission validates that submitted questions match expected questions
func ValidateQuestionSubmission(submitted dto.SolveQuestionRequest, expected entity.Question) error {
	if expected.ID != submitted.QuestionID {
		return fmt.Errorf("%w: %s", ErrValidation, invalidQuestionsError)
	}
	return nil
}

// ValidateGetQuizzesByUserAndTeamRequest validates the get quizzes by user and team request
func ValidateGetQuizzesByUserAndTeamRequest(userId string, teamId string, pageSize int) error {
	if userId == "" {
		return fmt.Errorf("%w: %s", ErrValidation, userIdEmptyError)
	}
	if teamId == "" {
		return fmt.Errorf("%w: %s", ErrValidation, teamIdEmptyError)
	}
	if pageSize <= 0 {
		return fmt.Errorf("%w: %s", ErrValidation, pageSizeInvalidError)
	}
	return nil
}

// ValidateGetQuizzesByTeamRequest validates the get quizzes by team request
func ValidateGetQuizzesByTeamRequest(userId string, teamId string, pageSize int) error {
	if userId == "" {
		return fmt.Errorf("%w: %s", ErrValidation, userIdEmptyError)
	}
	if teamId == "" {
		return fmt.Errorf("%w: %s", ErrValidation, teamIdEmptyError)
	}
	if pageSize <= 0 {
		return fmt.Errorf("%w: %s", ErrValidation, pageSizeInvalidError)
	}
	return nil
}
