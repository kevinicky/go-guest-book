package adapter

import (
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/kevinicky/go-guest-book/util"
	"sync"
)

type UserAdapter interface {
	CreateUser(req entity.CreateUserRequest) (*entity.UserSingleResponse, []error)
	GetUser(userID string) (*entity.UserSingleResponse, error)
	GetUsers(limit, offset int, key, isAdmin string) (*entity.UserMultiResponse, error)
	DeleteUser(userID string) error
	UpdateUser(userID string, req entity.UpdateUserRequest) (*entity.UserSingleResponse, []error)
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

func (u *userAdapter) GetUsers(limit, offset int, key, isAdmin string) (*entity.UserMultiResponse, error) {
	users, err := u.userUseCase.GetUsers(limit, offset, key, isAdmin)
	if err != nil {
		return nil, err
	}

	return u.setMultiUserResponse(users, limit, offset, key, isAdmin), nil
}

func (u *userAdapter) DeleteUser(userID string) error {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return err
	}

	return u.userUseCase.DeleteUser(userUUID)
}

func (u *userAdapter) UpdateUser(userID string, req entity.UpdateUserRequest) (*entity.UserSingleResponse, []error) {
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, []error{err}
	}

	user, errList := u.userUseCase.UpdateUser(userUUID, req)
	if errList != nil {
		return nil, errList
	}

	return u.setSingleUserResponse(user), nil
}

func (u *userAdapter) setSingleUserResponse(user *entity.User) *entity.UserSingleResponse {
	resp := entity.UserSingleResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	return &resp
}

func (u *userAdapter) setMultiUserResponse(users []entity.User, limit, offset int, key, isAdmin string) *entity.UserMultiResponse {
	var usersSingleResp []entity.UserSingleResponse
	var totalRows, totalPage, page int64
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, user := range users {
			usersSingleResp = append(usersSingleResp, *u.setSingleUserResponse(&user))
		}
	}()

	go func() {
		defer wg.Done()
		totalRows, _ = u.userUseCase.CountUser(key, isAdmin)
		totalPage, page = util.CountTotalPageAndCurrentPage(totalRows, limit, offset)
	}()

	wg.Wait()

	resp := entity.UserMultiResponse{
		Page:      page,
		Offset:    offset,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
		Filter: entity.UserQueryFilter{
			IsAdmin: isAdmin,
			Key:     key,
			OrderBy: entity.UserOrderBy{
				Field: "full_name",
				Sort:  "ascending",
			},
		},
		Rows: usersSingleResp,
	}

	return &resp
}
