package controller

import (
	"errors"
	"net/http"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
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
