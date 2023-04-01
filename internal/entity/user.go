package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UserSingleResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type UserMultiResponse struct {
	Page      int64                `json:"page"`
	Offset    int                  `json:"offset"`
	Limit     int                  `json:"limit"`
	TotalRows int64                `json:"total_rows"`
	TotalPage int64                `json:"total_page"`
	Filter    UserQueryFilter      `json:"filter"`
	Rows      []UserSingleResponse `json:"rows"`
}

type UserQueryFilter struct {
	Key     string      `json:"key"`
	IsAdmin string      `json:"is_admin"`
	OrderBy UserOrderBy `json:"order_by"`
}

type UserOrderBy struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}
