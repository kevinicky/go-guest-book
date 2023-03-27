package adapter

import (
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
)

type UserAdapter interface {
	CreateUser(req entity.CreateUserRequest) (*entity.UserSingleResponse, []error)
	GetUser(userID string) (*entity.UserSingleResponse, error)
}

type userAdapter struct {
	userUseCase usecase.UserUseCase
}

func NewUserAdapter(userUseCase usecase.UserUseCase) UserAdapter {
	return &userAdapter{
		userUseCase: userUseCase,
	}
}

func (u *userAdapter) CreateUser(req entity.CreateUserRequest) (*entity.UserSingleResponse, []error) {
	user, err := u.userUseCase.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return u.setSingleUserResponse(user), nil
}

func (u *userAdapter) GetUser(userID string) (*entity.UserSingleResponse, error) {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userUseCase.GetUser(userUUID)
	if err != nil {
		return nil, err
	}

	return u.setSingleUserResponse(user), nil
}

func (u *userAdapter) setSingleUserResponse(user *entity.User) *entity.UserSingleResponse {
	resp := entity.UserSingleResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	return &resp
}
