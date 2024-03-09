package services

import (
	"errors"
	"forum/Api/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (auth *AuthService) Authenticate(Password string, toAuthenticate *models.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(toAuthenticate.Password))
	if err != nil {
		return errors.New("invalide email or password")
	}
	return nil
}
