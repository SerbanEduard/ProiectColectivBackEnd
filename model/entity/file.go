package entity

type File struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Extension string `json:"extension"`
	Content   string `json:"content,omitempty"`
	Size      int64  `json:"size"`
	OwnerID   string `json:"ownerId,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

func NewFile(id, name, ftype, extension, content, ownerId string, size, createdAt, updatedAt int64) *File {
	return &File{
		ID:        id,
		Name:      name,
		Type:      ftype,
		Extension: extension,
		Content:   content,
		Size:      size,
		OwnerID:   ownerId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
