package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo UserRepositoryInterface
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: persistence.NewUserRepository(),
	}
}

func NewUserServiceWithRepo(repo interface{}) *UserService {
	return &UserService{
		userRepo: repo.(UserRepositoryInterface),
	}
}

type UserRepositoryInterface interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	GetAll() ([]*entity.User, error)
}

func (us *UserService) SignUp(request *dto.SignUpUserRequest) (*dto.SignUpUserResponse, error) {
	if err := validator.ValidateSignUpRequest(request); err != nil {
		return nil, err
	}

	if _, err := us.userRepo.GetByUsername(request.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	if _, err := us.userRepo.GetByEmail(request.Email); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(
		id,
		request.FirstName,
		request.LastName,
		request.Username,
		request.Email,
		string(hashedPassword),
		request.TopicsOfInterest,
	)

	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return dto.NewSignUpUserResponse(user.FirstName, user.LastName, user.Username), nil
}

// Login verifies credentials and returns basic user info on success.
func (us *UserService) Login(request *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	// Find user by email
	user, err := us.userRepo.GetByEmail(request.Email)
	if err != nil {
		// Do not leak whether email exists
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create JWT token
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is not configured")
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	resp := dto.NewLoginUserResponse(user.ID, user.FirstName, user.LastName, user.Username, user.TopicsOfInterest)
	resp.Token = signed
	return resp, nil
}

func (us *UserService) GetUserByID(id string) (*entity.User, error) {
	return us.userRepo.GetByID(id)
}

func (us *UserService) GetUserByEmail(email string) (*entity.User, error) {
	return us.userRepo.GetByEmail(email)
}

func (us *UserService) UpdateUser(user *entity.User) error {
	return us.userRepo.Update(user)
}

func (us *UserService) DeleteUser(id string) error {
	return us.userRepo.Delete(id)
}

func (us *UserService) GetAllUsers() ([]*entity.User, error) {
	return us.userRepo.GetAll()
}

func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
