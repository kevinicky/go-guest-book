package usecase

import "github.com/kevinicky/go-guest-book/internal/repository"

type GuestBookUseCase interface {
	ServerHealth() string
	DBHealth() string
}

type guestBookUseCase struct {
	guestBookRepository repository.GuestBookRepository
}

func NewGuestBookUseCase(guestBookRepository repository.GuestBookRepository) GuestBookUseCase {
	return &guestBookUseCase{
		guestBookRepository: guestBookRepository,
	}
}

func (gb *guestBookUseCase) ServerHealth() string {
	return "ok"
}

func (gb *guestBookUseCase) DBHealth() string {
	err := gb.guestBookRepository.Ping()
	if err == nil {
		return "ok"
	}

	return err.Error()
}
