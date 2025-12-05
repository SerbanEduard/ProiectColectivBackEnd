package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService FileServiceInterface
}

func NewFileController() *FileController {
	return &FileController{fileService: service.NewFileService()}
}

func NewFileControllerWithService(svc FileServiceInterface) *FileController {
	return &FileController{fileService: svc}
}

type FileServiceInterface interface {
	CreateFile(request *dto.FileUploadRequest, userID string) (*dto.FileUploadResponse, error)
	GetFileByID(id, userID string) (*entity.File, error)
	GetFilesByTeam(teamID, userID string, page, limit int) (*dto.FileListResponse, error)
	DeleteFile(id, userID string) error
}

// UploadFile
//
//	@Summary	Upload a file to a team (base64 content)
//	@Security	Bearer
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"Team ID"
//	@Param		request	body		dto.FileUploadRequest	true	"File upload request"
//	@Success	201		{object}	dto.FileUploadResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	403		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/teams/{id}/files [post]
func (fc *FileController) UploadFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teamID := c.Param("id")

	var req dto.FileUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set context from URL path
	req.ContextType = "team"
	req.ContextID = teamID

	resp, err := fc.fileService.CreateFile(&req, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not a member") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetFile
//
//	@Summary	Get file by id (with content)
//	@Security	Bearer
//	@Produce	json
//	@Param		id		path		string	true	"Team ID"
//	@Param		fileId	path		string	true	"File ID"
//	@Success	200		{object}	entity.File
//	@Failure	403		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/teams/{id}/files/{fileId} [get]
func (fc *FileController) GetFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fileID := c.Param("fileId")
	file, err := fc.fileService.GetFileByID(fileID, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "not a member") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, file)
}

// GetFilesByTeam
//
//	@Summary	Get all files for a team (metadata only, paginated)
//	@Security	Bearer
//	@Produce	json
//	@Param		id		path		string	true	"Team ID"
//	@Param		page	query		int		false	"Page number (default 1)"
//	@Param		limit	query		int		false	"Items per page (default 10, max 100)"
//	@Success	200		{object}	dto.FileListResponse
//	@Failure	403		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/teams/{id}/files [get]
func (fc *FileController) GetFilesByTeam(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teamID := c.Param("id")
	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	resp, err := fc.fileService.GetFilesByTeam(teamID, userID.(string), page, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not a member") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteFile
//
//	@Summary	Delete a file
//	@Security	Bearer
//	@Param		id		path		string	true	"Team ID"
//	@Param		fileId	path		string	true	"File ID"
//	@Success	200		{object}	map[string]string
//	@Failure	403		{object}	map[string]string
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/teams/{id}/files/{fileId} [delete]
func (fc *FileController) DeleteFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fileID := c.Param("fileId")
	if err := fc.fileService.DeleteFile(fileID, userID.(string)); err != nil {
		if strings.Contains(err.Error(), "not a member") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
