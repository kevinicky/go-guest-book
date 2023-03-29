package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"gorm.io/gorm"
)

func newHealthRepository(db *gorm.DB) repository.HealthRepository {
	return repository.NewHealthRepository(db)
}

func newHealthUseCase(healthRepository repository.HealthRepository) usecase.HealthUseCase {
	return usecase.NewHealthUseCase(healthRepository)
}

func newHealthAdapter(healthUseCase usecase.HealthUseCase) adapter.HealthAdapter {
	return adapter.NewHealthAdapter(healthUseCase)
}
