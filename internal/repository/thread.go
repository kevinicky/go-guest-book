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

type ThreadRepository interface {
	FindThread(id uuid.UUID) (*entity.Thread, error)
	GetThreads(limit, offset int, visitID uuid.UUID) ([]entity.Thread, error)
	CountThread(visitID uuid.UUID) (int64, error)
	SoftDeleteThread(threadID uuid.UUID) error
	CreateThread(thread entity.Thread) (*entity.Thread, error)
	SetCacheData(ctx context.Context, key string, value []byte) (string, error)
	DeleteCacheData(ctx context.Context, key string) (int64, error)
	GetCacheData(ctx context.Context, key string) (string, error)
}

type threadRepository struct {
	pgDB     *gorm.DB
	redisDB  *redis.Client
	redisTTL time.Duration
}

func NewThreadRepository(pgDB *gorm.DB, redisDB *redis.Client, redisTTL time.Duration) ThreadRepository {
	return &threadRepository{
		pgDB:     pgDB,
		redisDB:  redisDB,
		redisTTL: redisTTL,
	}
}

func (t *threadRepository) FindThread(id uuid.UUID) (*entity.Thread, error) {
	var thread entity.Thread
	thread.ID = id

	resp := t.pgDB.Joins("join visits on threads.visit_id = visits.id").Where("visits.deleted_at = ?", time.Time{}).Where("threads.deleted_at = ?", time.Time{}).First(&thread)
	if resp.Error != nil {
		if resp.Error.Error() == "record not found" {
			resp.Error = errors.New(customerror.THREAD_NOT_FOUND)
		}
	}

	return &thread, resp.Error
}

func (t *threadRepository) GetThreads(limit, offset int, visitID uuid.UUID) ([]entity.Thread, error) {
	var threads []entity.Thread
	resp := t.pgDB.Limit(limit).Offset(offset).Joins("join visits on threads.visit_id = visits.id").Where("visits.deleted_at = ?", time.Time{}).Where("visits.id = ?", visitID.String()).Where("threads.deleted_at = ?", time.Time{}).Where("visits.deleted_at = ?", time.Time{}).Order("threads.created_at ASC").Find(&threads)

	return threads, resp.Error
}

func (t *threadRepository) CountThread(visitID uuid.UUID) (int64, error) {
	var count int64
	count = -1
	resp := t.pgDB.Model(entity.Thread{}).Joins("join visits on threads.visit_id = visits.id").Where("visits.deleted_at = ?", time.Time{}).Where("visits.id = ?", visitID.String()).Where("threads.deleted_at = ?", time.Time{}).Where("visits.deleted_at = ?", time.Time{}).Count(&count)

	return count, resp.Error
}

func (t *threadRepository) SoftDeleteThread(threadID uuid.UUID) error {
	resp := t.pgDB.Model(entity.Thread{ID: threadID}).Update("deleted_at", time.Now())
	if resp.RowsAffected != 1 {
		return errors.New(customerror.VISIT_NOT_FOUND)
	}

	return resp.Error
}

func (t *threadRepository) CreateThread(thread entity.Thread) (*entity.Thread, error) {
	resp := t.pgDB.Create(&thread)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &thread, nil
}

func (t *threadRepository) SetCacheData(ctx context.Context, key string, value []byte) (string, error) {
	res, err := t.redisDB.Set(ctx, key, value, t.redisTTL*time.Second).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (t *threadRepository) DeleteCacheData(ctx context.Context, key string) (int64, error) {
	res, err := t.redisDB.Del(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	return res, nil
}

func (t *threadRepository) GetCacheData(ctx context.Context, key string) (string, error) {
	res, err := t.redisDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
