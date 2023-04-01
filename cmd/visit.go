package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

func newVisitRepository(dbPG *gorm.DB, dbRedis *redis.Client, redisTTL time.Duration) repository.VisitRepository {
	return repository.NewVisitRepository(dbPG, dbRedis, redisTTL)
}

func newVisitUseCase(visitRepository repository.VisitRepository, userUseCase usecase.UserUseCase) usecase.VisitUseCase {
	return usecase.NewVisitUseCase(visitRepository, userUseCase)
}

func newVisitAdapter(visitUseCase usecase.VisitUseCase, userAdapter adapter.UserAdapter) adapter.VisitAdapter {
	return adapter.NewVisitAdapter(visitUseCase, userAdapter)
}
