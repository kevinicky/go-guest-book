package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type JwtClaims struct {
	CredentialID string
	Expired      time.Duration
	Issuer       string
	SecretKey    []byte
}

type JwtRequest struct {
	Credential string
	Issuer     string `json:"issuer"`
}

type JwtResponse struct {
	Credential uuid.UUID `json:"credential"`
	Issuer     string    `json:"issuer"`
	Expired    int64     `json:"expired"`
	Jwt        string    `json:"jwt"`
}

type Credentials struct {
	Email string `json:"email"`
}

type Auth struct {
	Email     string `json:"email"`
	Jwt       string `json:"jwt"`
	SecretKey string `json:"secret_key"`
	Issuer    string `json:"issuer"`
	Expired   int64  `json:"expired"`
}

type UserMatrix struct {
	Id        uuid.UUID
	Endpoint  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserMatrixValidateRequest struct {
	UserID   uuid.UUID
	Endpoint string
	IsAdmin  bool
}
