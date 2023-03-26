package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"gorm.io/gorm"
)

func newUserRepository(db *gorm.DB) repository.UserRepository {
	return repository.NewUserRepository(db)
}

func newUserUseCase(UserRepository repository.UserRepository) usecase.UserUseCase {
	return usecase.NewUserUseCase(UserRepository)
}

func newUserAdapter(UserUseCase usecase.UserUseCase) adapter.UserAdapter {
	return adapter.NewUserAdapter(UserUseCase)
}
