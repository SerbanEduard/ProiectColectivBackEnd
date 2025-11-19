package controller

import (
	"errors"
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type QuizController struct {
	quizService service.QuizServiceInterface
}

func NewQuizController() *QuizController {
	return &QuizController{
		quizService: service.NewQuizService(),
	}
}

func NewQuizControllerWithService(quizService service.QuizServiceInterface) *QuizController {
	return &QuizController{
		quizService: quizService,
	}
}

// CreateQuiz
//
//	@Summary	Create a new quiz
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		entity.Quiz	true	"The create quiz request"
//	@Success	201		{object}	dto.CreateQuizResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	403		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/quizzes [post]
func (qc *QuizController) CreateQuiz(c *gin.Context) {
	var request entity.Quiz
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := qc.quizService.CreateQuiz(request)
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetQuizWithAnswers
//
//	@Summary	Get a quiz with answers
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The id for quiz"
//	@Success	200	{object}	entity.Quiz
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/quizzes/{id} [get]
func (qc *QuizController) GetQuizWithAnswers(c *gin.Context) {
	id := c.Param("id")
	quiz, err := qc.quizService.GetQuizWithAnswersById(id)

	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// GetQuizWithoutAnswers
//
//	@Summary	Get a quiz without answers
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The id for quiz"
//	@Success	200	{object}	dto.ReadQuizRequest
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/quizzes/{id}/test [get]
func (qc *QuizController) GetQuizWithoutAnswers(c *gin.Context) {
	id := c.Param("id")
	quiz, err := qc.quizService.GetQuizWithoutAnswersById(id)

	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// SolveQuiz
//
//	@Summary	Solve a quiz
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.SolveQuizRequest	true	"The solve quiz request"
//	@Success	200		{object}	dto.SolveQuizResponse
//	@Failure	403		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/quizzes/{id}/test [post]
func (qc *QuizController) SolveQuiz(c *gin.Context) {
	var request dto.SolveQuizRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claimsI, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsI.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	response, err := qc.quizService.SolveQuiz(request, userID)

	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, response)
}
