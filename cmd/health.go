package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func newHealthRepository(dbPostgres *gorm.DB, dbRedis *redis.Client) repository.HealthRepository {
	return repository.NewHealthRepository(dbPostgres, dbRedis)
}

func newHealthUseCase(healthRepository repository.HealthRepository) usecase.HealthUseCase {
	return usecase.NewHealthUseCase(healthRepository)
}

func newHealthAdapter(healthUseCase usecase.HealthUseCase) adapter.HealthAdapter {
	return adapter.NewHealthAdapter(healthUseCase)
}
