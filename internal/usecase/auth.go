package usecase

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"strings"
	"time"
)

type AuthUseCase interface {
	CreateJWT(payloadAuth entity.JwtRequest) (*entity.JwtResponse, error)
	ValidateJWT(tokenString string) (*entity.JwtResponse, error)
	CheckCredentials(credential, password string) error
}

type authUseCase struct {
	userUseCase UserUseCase
	jwtClaims   entity.JwtClaims
}

func NewAuthUseCase(userUseCase UserUseCase, jwtClaims entity.JwtClaims) AuthUseCase {
	return &authUseCase{
		userUseCase: userUseCase,
		jwtClaims:   jwtClaims,
	}
}

func (a *authUseCase) CreateJWT(payloadAuth entity.JwtRequest) (*entity.JwtResponse, error) {
	if payloadAuth.Issuer == "" {
		return nil, errors.New(customerror.ISSUER_MANDATORY)
	}

	claims := a.jwtClaims
	claims.Issuer = payloadAuth.Issuer
	claims.CredentialID = payloadAuth.Credential
	jwtExpiredTime := time.Now().Add(claims.Expired * time.Second).Unix()

	userID, err := a.getUserID(payloadAuth.Credential)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":           jwtExpiredTime,
		"authorized":    true,
		"credential_id": claims.CredentialID,
		"issuer":        claims.Issuer,
	})

	tokenString, err := token.SignedString(claims.SecretKey)
	jwtResponse := entity.JwtResponse{
		Credential: userID,
		Issuer:     claims.Issuer,
		Expired:    jwtExpiredTime,
		Jwt:        tokenString,
	}

	return &jwtResponse, err
}

func (a *authUseCase) ValidateJWT(tokenString string) (*entity.JwtResponse, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}

		return a.jwtClaims.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, errUserID := a.getUserID(claims["credential_id"].(string))
		if errUserID != nil {
			return nil, errUserID
		}

		jwtResponse := entity.JwtResponse{
			Credential: userID,
			Issuer:     claims["issuer"].(string),
			Expired:    int64(claims["exp"].(float64)),
			Jwt:        tokenString,
		}

		return &jwtResponse, nil
	} else {
		return nil, errors.New(customerror.INVALID_CREDENTIAL)
	}
}

func (a *authUseCase) getUserID(credential string) (uuid.UUID, error) {
	user, err := a.userUseCase.GetUser(uuid.Nil, credential)
	if err != nil {
		return uuid.Nil, errors.New(customerror.INVALID_CREDENTIAL)
	}

	return user.ID, nil
}

func (a *authUseCase) CheckCredentials(credential, password string) error {
	user, err := a.userUseCase.GetUser(uuid.Nil, credential)
	if err != nil {
		return errors.New(customerror.INVALID_CREDENTIAL)
	}

	sanitisePassword := util.HashSHA256(password)
	if user.Password != sanitisePassword {
		return errors.New(customerror.INVALID_CREDENTIAL)
	}

	return nil
}
