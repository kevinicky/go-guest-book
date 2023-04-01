package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

func newThreadRepository(dbPG *gorm.DB, dbRedis *redis.Client, redisTTL time.Duration) repository.ThreadRepository {
	return repository.NewThreadRepository(dbPG, dbRedis, redisTTL)
}

func newThreadUseCase(threadRepository repository.ThreadRepository, visitUseCase usecase.VisitUseCase, userUseCase usecase.UserUseCase) usecase.ThreadUseCase {
	return usecase.NewThreadUseCase(threadRepository, visitUseCase, userUseCase)
}

func newThreadAdapter(threadUseCase usecase.ThreadUseCase, visitAdapter adapter.VisitAdapter, userAdapter adapter.UserAdapter) adapter.ThreadAdapter {
	return adapter.NewThreadAdapter(threadUseCase, visitAdapter, userAdapter)
}
