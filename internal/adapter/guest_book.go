package adapter

import (
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
)

type GuestBookAdapter interface {
	Health() entity.Health
}

type guestBookAdapter struct {
	guestBookUseCase usecase.GuestBookUseCase
}

func NewGuestBookAdapter(guestBookUseCase usecase.GuestBookUseCase) GuestBookAdapter {
	return &guestBookAdapter{
		guestBookUseCase: guestBookUseCase,
	}
}

func (gb *guestBookAdapter) Health() entity.Health {
	serverStatus := gb.guestBookUseCase.ServerHealth()
	dbStatus := gb.guestBookUseCase.DBHealth()

	return entity.Health{
		Status: entity.HealthComponent{
			Server:   serverStatus,
			Database: dbStatus,
		},
	}
}
