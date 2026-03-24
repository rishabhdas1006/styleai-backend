package service

import (
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"
	"styleai-backend/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewUserService(repo *repository.UserRepository, secret string) *UserService {
	return &UserService{
		UserRepo:  repo,
		JWTSecret: secret,
	}
}

func (s *UserService) Register(name string, email string, password string) (*models.User, error) {
	existingUser, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, common.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.UserRepo.CreateUser(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", common.ErrEmailNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", common.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Role, s.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
