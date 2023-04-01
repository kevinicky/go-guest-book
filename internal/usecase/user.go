package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/util"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"net/mail"
	"strings"
	"sync"
	"time"
	"unicode"
)

type UserUseCase interface {
	CreateUser(req entity.CreateUserRequest) (*entity.User, []error)
	GetUser(userID uuid.UUID, email string) (*entity.User, error)
	GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error)
	CountUser(key, isAdmin string) (int64, error)
	DeleteUser(userID uuid.UUID) error
	UpdateUser(userID uuid.UUID, req entity.UpdateUserRequest) (*entity.User, []error)
	GetUserMatrix(endpoint string, isAdmin bool) ([]entity.UserMatrix, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (u *userUseCase) GetUser(userID uuid.UUID, email string) (*entity.User, error) {
	keyCache := "user~" + userID.String()
	ctx := context.Background()
	userCache := entity.User{}

	res, errCache := u.userRepository.GetCacheData(ctx, keyCache)
	if errCache != nil {
		user, err := u.userRepository.FindUser(userID, email)
		if err != nil {
			return nil, err
		}

		resMarshall, _ := json.Marshal(user)
		_, _ = u.userRepository.SetCacheData(ctx, keyCache, resMarshall)

		return user, err
	}

	_ = json.Unmarshal([]byte(res), &userCache)

	return &userCache, nil
}

func (u *userUseCase) DeleteUser(userID uuid.UUID) error {
	ctx := context.Background()
	keyCache := "user~" + userID.String()
	_, _ = u.userRepository.DeleteCacheData(ctx, keyCache)

	_, err := u.userRepository.FindUser(userID, "")
	if err != nil {
		return err
	}

	return u.userRepository.SoftDeleteUser(userID)
}

func (u *userUseCase) GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error) {
	if isAdmin != "" && isAdmin != "true" && isAdmin != "false" {
		return nil, errors.New(customerror.IS_ADMIN_WRONG_CONTENT)
	}

	return u.userRepository.GetUsers(limit, offset, key, isAdmin)
}

func (u *userUseCase) CountUser(key, isAdmin string) (int64, error) {
	return u.userRepository.CountUser(key, isAdmin)
}

func (u *userUseCase) CreateUser(req entity.CreateUserRequest) (*entity.User, []error) {
	var errList []error
	var sanitiseFullName, sanitiseEmail, sanitisePassword string
	var err error
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		sanitiseFullName = u.sanitiseFullName(req.FullName)
		err = u.validateFullName(sanitiseFullName)
		if err != nil {
			errList = append(errList, err)
		}
	}()

	go func() {
		defer wg.Done()
		sanitiseEmail = u.sanitiseEmail(req.Email)
		err = u.validateEmail(sanitiseEmail)
		if err != nil {
			errList = append(errList, err)
		}
	}()

	go func() {
		defer wg.Done()
		err = u.validatePassword(req.Password)
		if err != nil {
			errList = append(errList, err)

			return
		}

		sanitisePassword = u.sanitisePassword(req.Password)
	}()

	wg.Wait()
	if len(errList) > 0 {
		return nil, errList
	}

	newID, _ := uuid.NewV4()
	user := entity.User{
		ID:        newID,
		FullName:  sanitiseFullName,
		Email:     sanitiseEmail,
		Password:  sanitisePassword,
		IsAdmin:   req.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}

	userResp, err := u.userRepository.CreateUser(user)
	if err != nil {
		errList = append(errList, err)
	}

	keyCache := "user~" + userResp.ID.String()
	resMarshall, _ := json.Marshal(userResp)
	_, _ = u.userRepository.SetCacheData(ctx, keyCache, resMarshall)

	return userResp, errList
}

func (u *userUseCase) UpdateUser(userID uuid.UUID, req entity.UpdateUserRequest) (*entity.User, []error) {
	var errList []error
	ctx := context.Background()

	oldUser, err := u.GetUser(userID, "")
	if err != nil {
		return nil, []error{err}
	}

	var sanitiseFullName, sanitiseEmail string
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if req.FullName != "" {
			sanitiseFullName = u.sanitiseFullName(req.FullName)
		}
	}()

	go func() {
		defer wg.Done()
		if req.Email != "" {
			sanitiseEmail = u.sanitiseEmail(req.Email)
		}
	}()

	wg.Wait()

	edited := false
	if sanitiseFullName != oldUser.FullName {
		edited = true
		err = u.validateFullName(req.FullName)

		if err != nil {
			errList = append(errList, err)
		}
	}

	if sanitiseEmail != oldUser.Email {
		edited = true
		err = u.validateEmail(req.Email)

		if err != nil {
			errList = append(errList, err)
		}
	}

	if !edited {
		errList = append(errList, errors.New(customerror.USER_NO_RECORD_CHANGED))
	}

	if len(errList) > 0 {
		return nil, errList
	}

	keyCache := "user~" + userID.String()
	_, _ = u.userRepository.DeleteCacheData(ctx, keyCache)

	user := entity.User{
		ID:        userID,
		FullName:  sanitiseFullName,
		Email:     sanitiseEmail,
		Password:  oldUser.Password,
		UpdatedAt: time.Now(),
	}

	err = u.userRepository.UpdateUser(user)
	if err != nil {
		errList = append(errList, err)
	}

	keyCache = "user~" + userID.String()
	resMarshall, _ := json.Marshal(user)
	_, _ = u.userRepository.SetCacheData(ctx, keyCache, resMarshall)

	return &user, errList
}

func (u *userUseCase) GetUserMatrix(endpoint string, isAdmin bool) ([]entity.UserMatrix, error) {
	return u.userRepository.GetUserMatrix(endpoint, isAdmin)
}

func (u *userUseCase) sanitiseEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	return email
}

func (u *userUseCase) validateEmail(email string) error {
	if email == "" {
		return errors.New(customerror.EMAIL_MANDATORY)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New(customerror.INVALID_EMAIL)
	}

	totalEmail, err := u.userRepository.CountExistingEmail(email)
	if err != nil {
		return err
	}

	if totalEmail > 0 {
		return errors.New(customerror.EMAIL_TAKEN)
	}

	return nil
}

func (u *userUseCase) sanitisePassword(password string) string {
	return util.HashSHA256(password)
}

func (u *userUseCase) validatePassword(password string) error {
	if len(password) > 64 {
		return errors.New(customerror.PASSWORD_LEN_GT_LIMIT)
	}

	if password == "" {
		return errors.New(customerror.PASSWORD_MANDATORY)
	}

	return nil
}

func (u *userUseCase) validateFullName(fullname string) error {
	if fullname == "" {
		return errors.New(customerror.FULL_NAME_MANDATORY)
	}

	return nil
}

func (u *userUseCase) sanitiseFullName(fullname string) string {
	fullname = strings.TrimSpace(fullname)
	words := strings.Split(fullname, " ")
	sanitiseFullname := ""

	for _, word := range words {
		word = strings.ToLower(word)
		r := []rune(word)
		r[0] = unicode.ToUpper(r[0])

		if sanitiseFullname == "" {
			sanitiseFullname += string(r)
		} else {
			sanitiseFullname += " " + string(r)
		}
	}

	return sanitiseFullname
}
