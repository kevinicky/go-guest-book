package repository

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/customerror"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"gorm.io/gorm"
	"strings"
)

type UserRepository interface {
	FindUser(id uuid.UUID) (*entity.User, error)
	GetAllUser(limit, offset int, key string, isAdmin bool, isActive bool) ([]entity.User, error)
	CountUser(key string, isAdmin bool, isActive bool) (int64, error)
	SoftDeleteUser(userID uuid.UUID) error
	RemoveDeletedAt(userID uuid.UUID) error
	CreateUser(user entity.User) (*entity.User, error)
	CountExistingPhoneNumber(username string) (int64, error)
	CountExistingEmail(email string) (int64, error)
	UpdateUser(user entity.User) (*entity.User, error)
}

type userRepository struct {
	pgDB *gorm.DB
}

func NewUserRepository(pgDB *gorm.DB) UserRepository {
	return &userRepository{
		pgDB: pgDB,
	}
}

func (u *userRepository) FindUser(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	user.ID = id

	resp := u.pgDB.Unscoped().First(&user)
	if resp.Error != nil {
		if resp.Error.Error() == "record not found" {
			resp.Error = errors.New(customerror.USER_NOT_FOUND)
		}
	}

	return &user, resp.Error
}

func (u *userRepository) GetAllUser(limit, offset int, key string, isAdmin bool, isActive bool) ([]entity.User, error) {
	var users []entity.User

	chain := u.pgDB.Unscoped().Limit(limit).Offset(offset)

	if key != "" {
		keyLike := strings.ToUpper("%" + key + "%")
		chain = chain.Where(chain.Or("UPPER(fullname) LIKE ?", keyLike).Or("UPPER(email) LIKE ?", keyLike).Or("UPPER(phone_number) LIKE ?", keyLike))
	}

	switch isAdmin {
	case true:
		chain = chain.Where("is_admin = ?", true)
	case false:
		chain = chain.Where("is_admin = ?", false)
	}

	switch isActive {
	case true:
		chain = chain.Where("deleted_at IS NULL")
	case false:
		chain = chain.Where("deleted_at is NOT NULL")
	}

	resp := chain.Find(&users)

	return users, resp.Error
}

func (u *userRepository) CountUser(key string, isAdmin bool, isActive bool) (int64, error) {
	var count int64
	count = -1
	chain := u.pgDB.Unscoped().Model(&entity.User{})

	if key != "" {
		keyLike := strings.ToUpper("%" + key + "%")
		chain = chain.Where(chain.Or("UPPER(fullname) LIKE ?", keyLike).Or("UPPER(email) LIKE ?", keyLike).Or("UPPER(phone_number) LIKE ?", keyLike))
	}

	switch isAdmin {
	case true:
		chain = chain.Where("is_admin = ?", true)
	case false:
		chain = chain.Where("is_admin = ?", false)
	}

	switch isActive {
	case true:
		chain = chain.Where("deleted_at IS NULL")
	case false:
		chain = chain.Where("deleted_at is NOT NULL")
	}

	resp := chain.Count(&count)

	return count, resp.Error
}

func (u *userRepository) SoftDeleteUser(userID uuid.UUID) error {
	resp := u.pgDB.Delete(&entity.User{}, userID)

	return resp.Error
}

func (u *userRepository) RemoveDeletedAt(userID uuid.UUID) error {
	resp := u.pgDB.Unscoped().Model(&entity.User{}).Where("id", userID).Update("deleted_at", nil)

	return resp.Error
}

func (u *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	resp := u.pgDB.Create(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return &user, nil
}

func (u *userRepository) CountExistingPhoneNumber(phoneNumber string) (int64, error) {
	var total int64
	total = -1
	resp := u.pgDB.Unscoped().Model(entity.User{}).Where("phone_number", phoneNumber).Count(&total)

	return total, resp.Error
}

func (u *userRepository) CountExistingEmail(email string) (int64, error) {
	var total int64
	total = -1
	resp := u.pgDB.Unscoped().Model(entity.User{}).Where("email", email).Count(&total)

	return total, resp.Error
}

func (u *userRepository) UpdateUser(user entity.User) (*entity.User, error) {
	resp := u.pgDB.Save(&user)

	return &user, resp.Error
}
