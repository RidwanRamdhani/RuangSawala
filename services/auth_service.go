package services

import (
	"database/sql"
	"errors"

	"github.com/ruangsawala/backend/models"
	"github.com/ruangsawala/backend/repositories"
	"github.com/ruangsawala/backend/utils"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{UserRepo: repositories.NewUserRepository(db)}
}

func (s *AuthService) RegisterUser(input models.RegisterReq) error {
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

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
	}
	auth := &models.Auth{
		PasswordHash: hashedPassword,
		IsVerified:   false,
	}

	return s.UserRepo.CreateUser(user, auth)
}

func (s *AuthService) LoginUser(input models.LoginReq) (string, error) {
	user, auth, err := s.UserRepo.GetUserByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if !utils.CompareHashAndPasswordString(auth.PasswordHash, input.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, 24)
	if err != nil {
		return "", err
	}

	return token, nil
}
