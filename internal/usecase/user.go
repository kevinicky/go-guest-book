package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/customerror"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/mail"
	"strings"
	"time"
)

type UserUseCase interface {
	CreateUser(req entity.CreateUserRequest) (*entity.User, []error)
	GetUser(userID uuid.UUID) (*entity.User, error)
	GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error)
	CountUser(key, isAdmin string) (int64, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (u *userUseCase) GetUser(userID uuid.UUID) (*entity.User, error) {
	return u.userRepository.FindUser(userID)
}

func (u *userUseCase) GetUsers(limit, offset int, key, isAdmin string) ([]entity.User, error) {
	if isAdmin != "" && isAdmin != "true" && isAdmin != "false" {
		return nil, errors.New(customerror.IS_ADMIN_WRONG_CONTENT)
	}

	return u.userRepository.GetAllUser(limit, offset, key, isAdmin)
}

func (u *userUseCase) CountUser(key, isAdmin string) (int64, error) {
	return u.userRepository.CountUser(key, isAdmin)
}

func (u *userUseCase) CreateUser(req entity.CreateUserRequest) (*entity.User, []error) {
	var errList []error
	sanitiseFullName, err := u.validateFullName(req.FullName)
	if err != nil {
		errList = append(errList, err)
	}

	sanitiseEmail, err := u.validateEmail(req.Email)
	if err != nil {
		errList = append(errList, err)
	}

	sanitisePhoneNumber, err := u.validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		errList = append(errList, err)
	}

	sanitisePassword, err := u.validatePassword(req.Password)
	if err != nil {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		return nil, errList
	}

	newID, _ := uuid.NewV4()
	user := entity.User{
		ID:          newID,
		FullName:    sanitiseFullName,
		Email:       sanitiseEmail,
		PhoneNumber: sanitisePhoneNumber,
		Password:    sanitisePassword,
		IsAdmin:     req.IsAdmin,
		CreatedAt:   time.Now(),
	}

	userResp, err := u.userRepository.CreateUser(user)
	if err != nil {
		errList = append(errList, err)
	}

	return userResp, errList
}

func (u *userUseCase) validatePhoneNumber(phoneNumber string) (string, error) {
	if phoneNumber == "" {
		return "", errors.New(customerror.PHONE_NUMBER_MANDATORY)
	}

	totalPhoneNumber, err := u.userRepository.CountExistingPhoneNumber(phoneNumber)
	if err != nil {
		return "", err
	}

	if totalPhoneNumber > 0 {
		return "", errors.New(customerror.PHONE_NUMBER_TAKEN)
	}

	return phoneNumber, nil
}

func (u *userUseCase) validateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)

	if email == "" {
		return "", errors.New(customerror.EMAIL_MANDATORY)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.New(customerror.INVALID_EMAIL)
	}

	totalEmail, err := u.userRepository.CountExistingEmail(email)
	if err != nil {
		return "", err
	}

	if totalEmail > 0 {
		return "", errors.New(customerror.EMAIL_TAKEN)
	}

	return email, nil
}

func (u *userUseCase) validatePassword(password string) (string, error) {
	if len(password) > 64 {
		return "", errors.New(customerror.PASSWORD_LEN_GT_LIMIT)
	}

	if password == "" {
		return "", errors.New(customerror.PASSWORD_MANDATORY)
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	newPasswordByte := hash.Sum(nil)

	return fmt.Sprintf("%x", newPasswordByte), nil
}

func (u *userUseCase) validateFullName(fullname string) (string, error) {
	if fullname == "" {
		return "", errors.New(customerror.FULL_NAME_MANDATORY)
	}

	fullname = strings.TrimSpace(fullname)
	caser := cases.Title(language.Indonesian)
	fullname = caser.String(fullname)

	return fullname, nil
}
