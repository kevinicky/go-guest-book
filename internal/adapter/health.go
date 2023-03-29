package adapter

import (
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
)

type HealthAdapter interface {
	Health() entity.Health
}

type healthAdapter struct {
	healthUseCase usecase.HealthUseCase
}

func NewHealthAdapter(healthUseCase usecase.HealthUseCase) HealthAdapter {
	return &healthAdapter{
		healthUseCase: healthUseCase,
	}
}

func (gb *healthAdapter) Health() entity.Health {
	serverStatus := gb.healthUseCase.ServerHealth()
	dbPGStatus := gb.healthUseCase.DBPGHealth()

	return entity.Health{
		Status: entity.HealthComponent{
			Server: serverStatus,
			Database: []entity.HealthDatabase{{
				Name:   "postgresql",
				Status: dbPGStatus,
			},
			},
		},
	}
}
