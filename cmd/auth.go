package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
)

func newAuthUseCase(userUseCase usecase.UserUseCase, jwtClaims entity.JwtClaims) usecase.AuthUseCase {
	return usecase.NewAuthUseCase(userUseCase, jwtClaims)
}

func newAuthAdapter(authUseCase usecase.AuthUseCase) adapter.AuthAdapter {
	return adapter.NewAuthAdapter(authUseCase)
}
