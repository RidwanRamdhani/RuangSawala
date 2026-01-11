package services

import (
	"database/sql"
	"errors"

	"github.com/ruangsawala/backend/repositories"
	"github.com/ruangsawala/backend/utils"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{UserRepo: repositories.NewUserRepository(db)}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (s *AuthService) RegisterUser(input RegisterReq) error {
	emailTaken, err := s.UserRepo.IsEmailTaken(input.Email)
	if err != nil {
		return err
	}
	if emailTaken {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	return s.UserRepo.CreateUser(input.Username, input.Email, hashedPassword)
}
