package repository

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserRepository interface {
	FindUser(id uuid.UUID, email string) (*entity.User, error)
	GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error)
	CountUser(key, isAdmin string) (int64, error)
	SoftDeleteUser(userID uuid.UUID) error
	CreateUser(user entity.User) (*entity.User, error)
	CountExistingEmail(email string) (int64, error)
	UpdateUser(user entity.User) error
	GetUserMatrix(endpoint string, isAdmin bool) ([]entity.UserMatrix, error)
	SetCacheData(ctx context.Context, key string, value []byte) (string, error)
	DeleteCacheData(ctx context.Context, key string) (int64, error)
	GetCacheData(ctx context.Context, key string) (string, error)
}

type userRepository struct {
	pgDB     *gorm.DB
	redisDB  *redis.Client
	redisTTL time.Duration
}

func NewUserRepository(pgDB *gorm.DB, redisDB *redis.Client, redisTTL time.Duration) UserRepository {
	return &userRepository{
		pgDB:     pgDB,
		redisDB:  redisDB,
		redisTTL: redisTTL,
	}
}

func (u *userRepository) FindUser(id uuid.UUID, email string) (*entity.User, error) {
	var user entity.User
	user.ID = id

	chain := u.pgDB.Where("deleted_at = ?", time.Time{})
	if email != "" {
		chain = chain.Where("email = ?", email)
	}
	resp := chain.First(&user)
	if resp.Error != nil {
		if resp.Error.Error() == "record not found" {
			resp.Error = errors.New(customerror.USER_NOT_FOUND)
		}
	}

	return &user, resp.Error
}

func (u *userRepository) GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error) {
	var users []entity.User

	chain := u.pgDB.Limit(limit).Offset(offset).Where("deleted_at = ?", time.Time{})

	if key != "" {
		keyLike := strings.ToUpper("%" + key + "%")
		chain = chain.Where(chain.Or("UPPER(full_name) LIKE ?", keyLike).Or("UPPER(email) LIKE ?", keyLike).Or("UPPER(phone_number) LIKE ?", keyLike))
	}

	if isAdmin == "true" {
		chain = chain.Where("is_admin = ?", true)
	} else if isAdmin == "false" {
		chain = chain.Where("is_admin = ?", false)
	}

	resp := chain.Order("full_name ASC").Find(&users)

	return users, resp.Error
}

func (u *userRepository) CountUser(key, isAdmin string) (int64, error) {
	var count int64
	count = -1
	chain := u.pgDB.Model(entity.User{}).Where("deleted_at = ?", time.Time{})

	if key != "" {
		keyLike := strings.ToUpper("%" + key + "%")
		chain = chain.Where(chain.Or("UPPER(full_name) LIKE ?", keyLike).Or("UPPER(email) LIKE ?", keyLike).Or("UPPER(phone_number) LIKE ?", keyLike))
	}

	if isAdmin == "true" {
		chain = chain.Where("is_admin = ?", true)
	} else if isAdmin == "false" {
		chain = chain.Where("is_admin = ?", false)
	}

	resp := chain.Count(&count)

	return count, resp.Error
}

func (u *userRepository) SoftDeleteUser(userID uuid.UUID) error {
	resp := u.pgDB.Model(entity.User{ID: userID}).Update("deleted_at", time.Now())
	if resp.RowsAffected != 1 {
		return errors.New(customerror.USER_NOT_FOUND)
	}

	return resp.Error
}

func (u *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	resp := u.pgDB.Create(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &user, nil
}

func (u *userRepository) CountExistingEmail(email string) (int64, error) {
	var total int64
	total = -1
	resp := u.pgDB.Model(entity.User{}).Where("email", email).Count(&total)

	return total, resp.Error
}

func (u *userRepository) UpdateUser(user entity.User) error {
	resp := u.pgDB.Save(user)

	return resp.Error
}

func (u *userRepository) GetUserMatrix(endpoint string, isAdmin bool) ([]entity.UserMatrix, error) {
	var usersMatrix []entity.UserMatrix
	resp := u.pgDB.Where("endpoint = ?", endpoint).Where("is_admin = ?", isAdmin).Find(&usersMatrix)
	if resp.RowsAffected < 1 {
		resp.Error = errors.New(customerror.USER_MATRIX_NOT_FOUND)
	}

	return usersMatrix, resp.Error
}

func (u *userRepository) SetCacheData(ctx context.Context, key string, value []byte) (string, error) {
	res, err := u.redisDB.Set(ctx, key, value, u.redisTTL*time.Second).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (u *userRepository) DeleteCacheData(ctx context.Context, key string) (int64, error) {
	res, err := u.redisDB.Del(ctx, key).Result()
	if err != nil {
		return -1, err
	}

	return res, nil
}

func (u *userRepository) GetCacheData(ctx context.Context, key string) (string, error) {
	res, err := u.redisDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
