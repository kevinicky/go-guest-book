package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"gorm.io/gorm"
)

func newThreadRepository(db *gorm.DB) repository.ThreadRepository {
	return repository.NewThreadRepository(db)
}

func newThreadUseCase(threadRepository repository.ThreadRepository, visitUseCase usecase.VisitUseCase, userUseCase usecase.UserUseCase) usecase.ThreadUseCase {
	return usecase.NewThreadUseCase(threadRepository, visitUseCase, userUseCase)
}

func newThreadAdapter(threadUseCase usecase.ThreadUseCase, visitAdapter adapter.VisitAdapter, userAdapter adapter.UserAdapter) adapter.ThreadAdapter {
	return adapter.NewThreadAdapter(threadUseCase, visitAdapter, userAdapter)
}
