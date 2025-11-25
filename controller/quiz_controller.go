package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/utils"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
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
//	@Tags		quizzes
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
		if errors.Is(err, validator.ErrValidation) {
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
//	@Tags		quizzes
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
		if errors.Is(err, validator.ErrValidation) {
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
//	@Tags		quizzes
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"The id for quiz"
//	@Success	200	{object}	dto.ReadQuizResponse
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/quizzes/{id}/test [get]
func (qc *QuizController) GetQuizWithoutAnswers(c *gin.Context) {
	id := c.Param("id")
	quiz, err := qc.quizService.GetQuizWithoutAnswersById(id)

	if err != nil {
		if errors.Is(err, validator.ErrValidation) {
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
//	@Tags		quizzes
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.SolveQuizRequest	true	"The solve quiz request"
//	@Param		id		path		string					true	"The id for quiz"
//	@Success	200		{object}	dto.SolveQuizResponse
//	@Failure	403		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/quizzes/{id}/test [post]
func (qc *QuizController) SolveQuiz(c *gin.Context) {
	quizID := c.Param("id")
	var request dto.SolveQuizRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	response, err := qc.quizService.SolveQuiz(request, userID, quizID)

	if err != nil {
		if errors.Is(err, validator.ErrValidation) {
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

// GetQuizzesByUserAndTeam
//
//	@Summary	Get quizzes by user with pagination
//	@Tags		quizzes
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		userId		path		string	true	"User ID"
//	@Param		teamId		path		string	true	"Team ID"
//	@Param		pageSize	query		int		false	"Page size (default 10)"
//	@Param		lastKey		query		string	false	"Last key for pagination"
//	@Success	200			{object}	map[string]interface{}
//	@Failure	400			{object}	map[string]string
//	@Failure	401			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/quizzes/user/{userId}/team/{teamId} [get]
func (qc *QuizController) GetQuizzesByUserAndTeam(c *gin.Context) {
	userID := c.Param("userId")
	teamID := c.Param("teamId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	lastKey := c.Query("lastKey")

	quizzes, newKey, err := qc.quizService.GetQuizzesByUserAndTeam(userID, teamID, pageSize, lastKey)
	if err != nil {
		if errors.Is(err, validator.ErrValidation) {
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

	response := map[string]interface{}{
		"quizzes": quizzes,
		"nextKey": newKey,
	}

	c.JSON(http.StatusOK, response)
}

// GetQuizzesByTeam
//
//	@Summary	Get quizzes by team with pagination
//	@Tags		quizzes
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		teamId		path		string	true	"Team ID"
//	@Param		pageSize	query		int		false	"Page size (default 10)"
//	@Param		lastKey		query		string	false	"Last key for pagination"
//	@Success	200			{object}	map[string]interface{}
//	@Failure	400			{object}	map[string]string
//	@Failure	401			{object}	map[string]string
//	@Failure	403			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/quizzes/team/{teamId} [get]
func (qc *QuizController) GetQuizzesByTeam(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team ID is required"})
		return
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	lastKey := c.Query("lastKey")

	quizzes, newKey, err := qc.quizService.GetQuizzesByTeam(userID, teamID, pageSize, lastKey)
	if err != nil {
		if errors.Is(err, validator.ErrValidation) {
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

	response := map[string]interface{}{
		"quizzes": quizzes,
		"nextKey": newKey,
	}

	c.JSON(http.StatusOK, response)
}
