package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type Visit struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateVisitRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type VisitSingleResponse struct {
	ID        uuid.UUID          `json:"id"`
	User      UserSingleResponse `json:"user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at"`
}

type VisitMultiResponse struct {
	Page      int64                 `json:"page"`
	Offset    int                   `json:"offset"`
	Limit     int                   `json:"limit"`
	TotalRows int64                 `json:"total_rows"`
	TotalPage int64                 `json:"total_page"`
	Filter    VisitQueryFilter      `json:"filter"`
	Rows      []VisitSingleResponse `json:"rows"`
}

type VisitQueryFilter struct {
	OrderBy VisitOrderBy `json:"order_by"`
}

type VisitOrderBy struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}
