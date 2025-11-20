package persistence

import (
	"context"
	"errors"
	"log"

	"firebase.google.com/go/v4/db"
	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const (
	quizCollection    = "quizzes"
	quizNotFoundError = "quiz not found"
	userIdField       = "user_id"
	teamIdField       = "team_id"
	userTeamIdField   = "user_team_id"
)

type QuizRepositoryInterface interface {
	Create(quiz entity.Quiz) error
	Update(quiz entity.Quiz) error
	GetById(id string) (entity.Quiz, error)
	GetByUser(id string, pageSize int, lastKey string) ([]entity.Quiz, string, error)
	GetByUserAndTeam(userId string, teamId string, pageSize int, lastKey string) ([]entity.Quiz, string, error)
	GetByTeam(id string, pageSize int, lastKey string) ([]entity.Quiz, string, error)
}

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository {
	return &QuizRepository{}
}

func (qr *QuizRepository) Create(quiz entity.Quiz) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection + "/" + quiz.ID)
	return ref.Set(ctx, quiz)
}

func (qr *QuizRepository) Update(quiz entity.Quiz) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection + "/" + quiz.ID)

	return ref.Set(ctx, quiz)
}

func (qr *QuizRepository) GetById(id string) (entity.Quiz, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection + "/" + id)

	var quiz entity.Quiz
	if err := ref.Get(ctx, &quiz); err != nil {
		return entity.Quiz{}, err
	}
	if quiz.ID == "" {
		return entity.Quiz{}, errors.New(quizNotFoundError)
	}

	return quiz, nil
}

func (qr *QuizRepository) GetByUser(id string, pageSize int, lastKey string) ([]entity.Quiz, string, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection)

	query := ref.OrderByChild(userIdField).EqualTo(id)

	if lastKey != "" {
		query = query.StartAt(lastKey)
	}

	limit := pageSize
	if lastKey != "" {
		limit = pageSize + 1
	}
	query = query.LimitToFirst(limit)

	return FilterByQuery(query, ctx, lastKey != "")
}

func (qr *QuizRepository) GetByTeam(id string, pageSize int, lastKey string) ([]entity.Quiz, string, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection)

	query := ref.OrderByChild(teamIdField).EqualTo(id)

	if lastKey != "" {
		query = query.StartAt(lastKey)
	}

	limit := pageSize
	if lastKey != "" {
		limit = pageSize + 1
	}
	query = query.LimitToFirst(limit)

	return FilterByQuery(query, ctx, lastKey != "")
}

func (qr *QuizRepository) GetByUserAndTeam(userId string, teamId string, pageSize int, lastKey string) ([]entity.Quiz, string, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(quizCollection)
	combinedId := userId + "_" + teamId

	query := ref.OrderByChild(userTeamIdField).EqualTo(combinedId)

	if lastKey != "" {
		query = query.StartAt(lastKey)
	}

	limit := pageSize
	if lastKey != "" {
		limit = pageSize + 1
	}
	query = query.LimitToFirst(limit)

	return FilterByQuery(query, ctx, lastKey != "")
}

func FilterByQuery(query *db.Query, ctx context.Context, hasCursor bool) ([]entity.Quiz, string, error) {
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, "", err
	}

	if len(results) == 0 {
		return []entity.Quiz{}, "", nil
	}

	startIndex := 0
	if hasCursor {
		if len(results) > 0 {
			startIndex = 1
		}
	}

	quizzes := make([]entity.Quiz, 0, len(results)-startIndex)
	for i := startIndex; i < len(results); i++ {
		result := results[i]

		var quiz entity.Quiz
		if err := result.Unmarshal(&quiz); err != nil {
			log.Printf("FilterByQuery: Error unmarshalling quiz at index %d: %v", i, err)
			return nil, "", err
		}

		quizzes = append(quizzes, quiz)
	}

	var newLastKey string
	if len(quizzes) > 0 {
		newLastKey = quizzes[len(quizzes)-1].ID
	}

	return quizzes, newLastKey, nil
}
