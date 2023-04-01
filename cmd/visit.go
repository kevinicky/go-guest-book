package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"gorm.io/gorm"
)

func newVisitRepository(db *gorm.DB) repository.VisitRepository {
	return repository.NewVisitRepository(db)
}

func newVisitUseCase(visitRepository repository.VisitRepository, userUseCase usecase.UserUseCase) usecase.VisitUseCase {
	return usecase.NewVisitUseCase(visitRepository, userUseCase)
}

func newVisitAdapter(visitUseCase usecase.VisitUseCase, userAdapter adapter.UserAdapter) adapter.VisitAdapter {
	return adapter.NewVisitAdapter(visitUseCase, userAdapter)
}
