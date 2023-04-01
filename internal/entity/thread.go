package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type Thread struct {
	ID        uuid.UUID `json:"id"`
	VisitID   uuid.UUID `json:"visit_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateThreadRequest struct {
	VisitID uuid.UUID `json:"visit_id"`
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
}

type ThreadSingleResponse struct {
	ID        uuid.UUID           `json:"id"`
	Visit     VisitSingleResponse `json:"visit"`
	User      UserSingleResponse  `json:"user"`
	Content   string              `json:"content"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt time.Time           `json:"deleted_at"`
}

type ThreadMultiResponsePerRow struct {
	ID        uuid.UUID          `json:"id"`
	User      UserSingleResponse `json:"user"`
	Content   string             `json:"content"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at"`
}

type ThreadMultiResponse struct {
	Page      int64                       `json:"page"`
	Offset    int                         `json:"offset"`
	Limit     int                         `json:"limit"`
	TotalRows int64                       `json:"total_rows"`
	TotalPage int64                       `json:"total_page"`
	Filter    ThreadQueryFilter           `json:"filter"`
	Rows      []ThreadMultiResponsePerRow `json:"rows"`
}

type ThreadQueryFilter struct {
	VisitID uuid.UUID     `json:"visit_id"`
	OrderBy ThreadOrderBy `json:"order_by"`
}

type ThreadOrderBy struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}
