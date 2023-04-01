package adapter

import (
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"sync"
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
	var serverStatus, dbPGStatus, dbRedisStatus string
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		serverStatus = gb.healthUseCase.ServerHealth()
	}()
	go func() {
		defer wg.Done()
		dbPGStatus = gb.healthUseCase.DBPGHealth()
	}()
	go func() {
		defer wg.Done()
		dbRedisStatus = gb.healthUseCase.DBRedisHealth()
	}()
	
	wg.Wait()

	return entity.Health{
		Status: entity.HealthComponent{
			Server: serverStatus,
			Database: []entity.HealthDatabase{{
				Name:   "postgresql",
				Status: dbPGStatus,
			}, {
				Name:   "redis",
				Status: dbRedisStatus,
			},
			},
		},
	}
}
