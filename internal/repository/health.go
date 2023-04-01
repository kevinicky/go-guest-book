package repository

import "gorm.io/gorm"

type HealthRepository interface {
	PingPG() error
}

type healthRepository struct {
	dbPG *gorm.DB
}

func NewHealthRepository(dbPG *gorm.DB) HealthRepository {
	return &healthRepository{
		dbPG: dbPG,
	}
}

func (h *healthRepository) PingPG() error {
	db, err := h.dbPG.DB()
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
