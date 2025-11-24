package persistence

import (
	"context"
	"errors"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

const (
	filesCollection = "files"
	fileNotFound    = "file not found"
)

type FileRepository struct{}

type FileRepositoryInterface interface {
	Create(file *entity.File) error
	GetByID(id string) (*entity.File, error)
	GetAll() ([]*entity.File, error)
	Update(file *entity.File) error
	Delete(id string) error
}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (fr *FileRepository) Create(file *entity.File) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(filesCollection + "/" + file.ID)
	return ref.Set(ctx, file)
}

func (fr *FileRepository) GetByID(id string) (*entity.File, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(filesCollection + "/" + id)

	var file entity.File
	if err := ref.Get(ctx, &file); err != nil {
		return nil, err
	}
	if file.ID == "" {
		return nil, errors.New(fileNotFound)
	}
	return &file, nil
}

func (fr *FileRepository) GetAll() ([]*entity.File, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(filesCollection)

	var filesMap map[string]*entity.File
	if err := ref.Get(ctx, &filesMap); err != nil {
		return nil, err
	}

	files := make([]*entity.File, 0, len(filesMap))
	for _, f := range filesMap {
		files = append(files, f)
	}
	return files, nil
}

func (fr *FileRepository) Update(file *entity.File) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(filesCollection + "/" + file.ID)
	return ref.Set(ctx, file)
}

func (fr *FileRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef(filesCollection + "/" + id)
	return ref.Delete(ctx)
}
