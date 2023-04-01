package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthRepository interface {
	PingPG() error
	PingRedis() error
}

type healthRepository struct {
	dbPG    *gorm.DB
	dbRedis *redis.Client
}

func NewHealthRepository(dbPG *gorm.DB, dbRedis *redis.Client) HealthRepository {
	return &healthRepository{
		dbPG:    dbPG,
		dbRedis: dbRedis,
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

func (h *healthRepository) PingRedis() error {
	res := h.dbRedis.Ping(context.Background())
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
