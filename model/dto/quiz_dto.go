package dto

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type CreateQuizResponse struct {
	QuizID string `json:"quiz_id"`
}

type SolveQuestionRequest struct {
	QuestionID string   `json:"quiz_question_id"`
	Answer     []string `json:"answer"`
}

type SolveQuestionResponse struct {
	QuestionID    string   `json:"quiz_question_id"`
	IsCorrect     bool     `json:"is_correct"`
	CorrectFields []string `json:"correct_fields"`
}

type SolveQuizRequest struct {
	QuizID   string                 `json:"quiz_id"`
	Attempts []SolveQuestionRequest `json:"attempts"`
}

type SolveQuizResponse struct {
	IsCorrect         bool                    `json:"is_correct"`
	QuestionResponses []SolveQuestionResponse `json:"questions_answers"`
}

type ReadQuizRequest struct {
	QuizID string `json:"quiz_id"`
}

type ReadQuizQuestionResponse struct {
	QuestionID string   `json:"quiz_question_id"`
	Question   string   `json:"question"`
	Options    []string `json:"quiz_options"`
}

type ReadQuizResponse struct {
	QuizID        string                     `json:"quiz_id"`
	QuizTitle     string                     `json:"quiz_title"`
	QuizType      model.QuizType             `json:"quiz_type"`
	QuizQuestions []ReadQuizQuestionResponse `json:"quiz_questions"`
}

func NewSolveQuestionResponse(questionID string, isCorrect bool, correctFields []string) SolveQuestionResponse {
	return SolveQuestionResponse{
		QuestionID:    questionID,
		IsCorrect:     isCorrect,
		CorrectFields: correctFields,
	}
}

func NewReadQuizQuestionResponse(questionID string, question string, options []string) ReadQuizQuestionResponse {
	return ReadQuizQuestionResponse{
		QuestionID: questionID,
		Question:   question,
		Options:    options,
	}
}

func NewReadQuizResponse(quizID string, quizTitle string, quizQuestions []ReadQuizQuestionResponse) ReadQuizResponse {
	return ReadQuizResponse{
		QuizID:        quizID,
		QuizTitle:     quizTitle,
		QuizQuestions: quizQuestions,
	}
}

func NewCreateQuizResponse(quizID string) CreateQuizResponse {
	return CreateQuizResponse{QuizID: quizID}
}
