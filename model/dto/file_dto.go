package dto

type FileUploadRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Extension   string `json:"extension" binding:"required"`
	Content     string `json:"content" binding:"required"`
	OwnerID     string `json:"ownerId" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
	ContextType string `json:"contextType"` // Set automatically from URL ("team" or "chat")
	ContextID   string `json:"contextId"`   // Set automatically from URL (teamId or chatId)
}

type FileUploadResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Extension   string `json:"extension"`
	Size        int64  `json:"size"`
	OwnerID     string `json:"ownerId"`
	ContextType string `json:"contextType"`
	ContextID   string `json:"contextId"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type FileListResponse struct {
	Files      []*FileUploadResponse `json:"files"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalCount int                   `json:"totalCount"`
	TotalPages int                   `json:"totalPages"`
}
