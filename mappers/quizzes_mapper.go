package mappers

import (
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

func MapDomainToReadDTO(quiz entity.Quiz) dto.ReadQuizResponse {
	questions := make([]dto.ReadQuizQuestionResponse, len(quiz.Questions))
	for i := range quiz.Questions {
		questions[i] = dto.NewReadQuizQuestionResponse(
			quiz.Questions[i].ID,
			quiz.Questions[i].Question,
			quiz.Questions[i].Options,
		)
	}
	return dto.ReadQuizResponse{
		QuizID:        quiz.ID,
		QuizTitle:     quiz.QuizName,
		QuizQuestions: questions,
	}
}
