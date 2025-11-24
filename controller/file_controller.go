package controller

import (
	"net/http"
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
	CreateFile(request *dto.FileUploadRequest) (*dto.FileResponse, error)
	GetFileByID(id string) (*entity.File, error)
	GetAllFiles() ([]*entity.File, error)
	DeleteFile(id string) error
}

// UploadFile
//
//	@Summary	Upload a file (base64 content)
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dto.FileUploadRequest	true	"File upload request"
//	@Success	201		{object}	dto.FileResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/files [post]
func (fc *FileController) UploadFile(c *gin.Context) {
	var req dto.FileUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := fc.fileService.CreateFile(&req)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetFile
//
//	@Summary	Get file by id
//	@Produce	json
//	@Param		id	path		string	true	"File ID"
//	@Success	200	{object}	entity.File
//	@Failure	404	{object}	map[string]string
//	@Router		/files/{id} [get]
func (fc *FileController) GetFile(c *gin.Context) {
	id := c.Param("id")
	file, err := fc.fileService.GetFileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, file)
}

// GetAllFiles
//
//	@Summary	Get all files
//	@Produce	json
//	@Success	200	{array}	entity.File
//	@Router		/files [get]
func (fc *FileController) GetAllFiles(c *gin.Context) {
	files, err := fc.fileService.GetAllFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

// DeleteFile
//
//	@Summary	Delete a file
//	@Param		id	path		string	true	"File ID"
//	@Success	200	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Router		/files/{id} [delete]
func (fc *FileController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if err := fc.fileService.DeleteFile(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
