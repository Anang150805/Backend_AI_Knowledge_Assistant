package services

import (
	"backend-AI-Knowledge-Assistant/internal/models"
	"backend-AI-Knowledge-Assistant/internal/repositories"
	"backend-AI-Knowledge-Assistant/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) Register(username, password string) (*models.User, string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.UserRepo.Create(username, string(hashed))
	if err != nil {
		return nil, "", err
	}

	jwtToken, err := token.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", err
	}

	jwtToken, err := token.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}