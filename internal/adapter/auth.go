package adapter

import (
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
)

type AuthAdapter interface {
	CreateJWT(payloadAuth entity.JwtRequest) (*entity.JwtResponse, error)
	ValidateJWT(tokenString string) (*entity.JwtResponse, error)
	CheckCredentials(credential, password string) error
}

type authAdapter struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthAdapter(authUseCase usecase.AuthUseCase) AuthAdapter {
	return &authAdapter{authUseCase: authUseCase}
}

func (a *authAdapter) CreateJWT(payloadAuth entity.JwtRequest) (*entity.JwtResponse, error) {
	return a.authUseCase.CreateJWT(payloadAuth)
}

func (a *authAdapter) ValidateJWT(tokenString string) (*entity.JwtResponse, error) {
	return a.authUseCase.ValidateJWT(tokenString)
}

func (a *authAdapter) CheckCredentials(credential, password string) error {
	return a.authUseCase.CheckCredentials(credential, password)
}
