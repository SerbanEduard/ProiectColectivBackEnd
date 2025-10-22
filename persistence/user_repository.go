package persistence

import (
	"context"
	"fmt"

	"github.com/SerbanEduard/ProiectColectivBackEnd/config"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Create(user *entity.User) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users/" + user.ID)
	return ref.Set(ctx, user)
}

func (ur *UserRepository) GetByID(id string) (*entity.User, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users/" + id)
	
	var user entity.User
	if err := ref.Get(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users")
	
	query := ref.OrderByChild("email").EqualTo(email)
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	
	var user entity.User
	if err := results[0].Unmarshal(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Update(user *entity.User) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users/" + user.ID)
	return ref.Set(ctx, user)
}

func (ur *UserRepository) Delete(id string) error {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users/" + id)
	return ref.Delete(ctx)
}

func (ur *UserRepository) GetByUsername(username string) (*entity.User, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users")
	
	query := ref.OrderByChild("username").EqualTo(username)
	results, err := query.GetOrdered(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	
	var user entity.User
	if err := results[0].Unmarshal(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetAll() ([]*entity.User, error) {
	ctx := context.Background()
	ref := config.FirebaseDB.NewRef("users")
	
	var usersMap map[string]*entity.User
	if err := ref.Get(ctx, &usersMap); err != nil {
		return nil, err
	}
	
	users := make([]*entity.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, user)
	}
	return users, nil
}