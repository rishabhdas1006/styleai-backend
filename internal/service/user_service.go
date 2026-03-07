package service

import (
	"errors"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (s *UserService) Register(name string, email string, password string) (*models.User, error) {
	existingUser, _ := s.UserRepo.GetUserByEmail(email)

	if existingUser != nil {
		return nil, errors.New("email already exists")
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

func (s *UserService) Login(email string, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
