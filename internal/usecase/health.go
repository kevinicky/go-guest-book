package usecase

import "github.com/kevinicky/go-guest-book/internal/repository"

type HealthUseCase interface {
	ServerHealth() string
	DBPGHealth() string
}

type healthUseCase struct {
	healthRepository repository.HealthRepository
}

func NewHealthUseCase(healthRepository repository.HealthRepository) HealthUseCase {
	return &healthUseCase{
		healthRepository: healthRepository,
	}
}

func (h *healthUseCase) ServerHealth() string {
	return "ok"
}

func (h *healthUseCase) DBPGHealth() string {
	err := h.healthRepository.PingPG()
	if err == nil {
		return "ok"
	}

	return err.Error()
}
