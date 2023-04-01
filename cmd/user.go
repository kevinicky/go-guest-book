package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

func newUserRepository(dbPG *gorm.DB, dbRedis *redis.Client, redisTTL time.Duration) repository.UserRepository {
	return repository.NewUserRepository(dbPG, dbRedis, redisTTL)
}

func newUserUseCase(userRepository repository.UserRepository) usecase.UserUseCase {
	return usecase.NewUserUseCase(userRepository)
}

func newUserAdapter(userUseCase usecase.UserUseCase) adapter.UserAdapter {
	return adapter.NewUserAdapter(userUseCase)
}
