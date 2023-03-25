package main

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"gorm.io/gorm"
)

func newGuestBookRepository(db *gorm.DB) repository.GuestBookRepository {
	return repository.NewGuestBookRepository(db)
}

func newGuestBookUseCase(guestBookRepository repository.GuestBookRepository) usecase.GuestBookUseCase {
	return usecase.NewGuestBookUseCase(guestBookRepository)
}

func newGuestBookAdapter(guestBookUseCase usecase.GuestBookUseCase) adapter.GuestBookAdapter {
	return adapter.NewGuestBookAdapter(guestBookUseCase)
}
