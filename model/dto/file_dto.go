package dto

type FileUploadRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Content   string `json:"content" binding:"required"`
	OwnerID   string `json:"ownerId" binding:"required"`
	Size      int64  `json:"size" binding:"required"`
}

type FileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	OwnerID   string `json:"ownerId"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type FileListResponse struct {
	Files []*FileResponse `json:"files"`
}
