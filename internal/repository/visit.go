package repository

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type VisitRepository interface {
	FindVisit(id uuid.UUID) (*entity.Visit, error)
	GetVisits(limit, offset int) ([]entity.Visit, error)
	CountVisit() (int64, error)
	SoftDeleteVisit(visitID uuid.UUID) error
	CreateVisit(visit entity.Visit) (*entity.Visit, error)
	SetCacheData(ctx context.Context, key string, value []byte) (string, error)
	DeleteCacheData(ctx context.Context, key string) (int64, error)
	GetCacheData(ctx context.Context, key string) (string, error)
}

type visitRepository struct {
	pgDB     *gorm.DB
	redisDB  *redis.Client
	redisTTL time.Duration
}

func NewVisitRepository(pgDB *gorm.DB, redisDB *redis.Client, redisTTL time.Duration) VisitRepository {
	return &visitRepository{
		pgDB:     pgDB,
		redisDB:  redisDB,
		redisTTL: redisTTL,
	}
}

func (v *visitRepository) FindVisit(id uuid.UUID) (*entity.Visit, error) {
	var visit entity.Visit
	visit.ID = id

	resp := v.pgDB.Joins("join users on visits.user_id = users.id").Where("visits.deleted_at = ?", time.Time{}).Where("users.deleted_at = ?", time.Time{}).First(&visit)
	if resp.Error != nil {
		if resp.Error.Error() == "record not found" {
			resp.Error = errors.New(customerror.VISIT_NOT_FOUND)
		}
	}

	return &visit, resp.Error
}

func (v *visitRepository) GetVisits(limit, offset int) ([]entity.Visit, error) {
	var visits []entity.Visit
	resp := v.pgDB.Limit(limit).Offset(offset).Joins("join users on visits.user_id = users.id").Where("visits.deleted_at = ?", time.Time{}).Where("users.deleted_at = ?", time.Time{}).Order("visits.created_at ASC").Find(&visits)

	return visits, resp.Error
}

func (v *visitRepository) CountVisit() (int64, error) {
	var count int64
	count = -1
	resp := v.pgDB.Model(entity.Visit{}).Joins("join users on visits.user_id = users.id").Where("visits.deleted_at = ?", time.Time{}).Where("users.deleted_at = ?", time.Time{}).Count(&count)

	return count, resp.Error
}

func (v *visitRepository) SoftDeleteVisit(visitID uuid.UUID) error {
	resp := v.pgDB.Model(entity.Visit{ID: visitID}).Update("deleted_at", time.Now())
	if resp.RowsAffected != 1 {
		return errors.New(customerror.VISIT_NOT_FOUND)
	}

	return resp.Error
}

func (v *visitRepository) CreateVisit(visit entity.Visit) (*entity.Visit, error) {
	resp := v.pgDB.Create(&visit)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &visit, nil
}

func (v *visitRepository) SetCacheData(ctx context.Context, key string, value []byte) (string, error) {
	res, err := v.redisDB.Set(ctx, key, value, v.redisTTL*time.Second).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (v *visitRepository) DeleteCacheData(ctx context.Context, key string) (int64, error) {
	res, err := v.redisDB.Del(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	return res, nil
}

func (v *visitRepository) GetCacheData(ctx context.Context, key string) (string, error) {
	res, err := v.redisDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
