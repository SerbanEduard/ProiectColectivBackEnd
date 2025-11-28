package entity

import "github.com/SerbanEduard/ProiectColectivBackEnd/model"

type Question struct {
	ID       string         `json:"id,omitempty"`
	Type     model.QuizType `json:"type"`
	Question string         `json:"question"`
	Answers  []string       `json:"answers"`
	Options  []string       `json:"options"`
}

type Quiz struct {
	ID         string     `json:"id"`
	QuizName   string     `json:"quiz_name"`
	UserID     string     `json:"user_id"`
	TeamID     string     `json:"team_id"`
	Questions  []Question `json:"questions"`
	UserTeamId string     `json:"user_team_id"`
}

func NewQuestion(ID string, quizType model.QuizType, question string, answers []string, options []string) *Question {
	return &Question{
		Type:     quizType,
		Question: question,
		Answers:  answers,
		Options:  options,
	}
}

func NewQuiz(ID string, quizName string, userID string, teamID string, questions []Question) *Quiz {
	return &Quiz{
		ID:         ID,
		QuizName:   quizName,
		UserID:     userID,
		TeamID:     teamID,
		Questions:  questions,
		UserTeamId: userID + "_" + teamID,
	}
}
