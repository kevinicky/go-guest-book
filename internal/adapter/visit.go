package adapter

import (
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/kevinicky/go-guest-book/util"
	"sync"
)

type VisitAdapter interface {
	CreateVisit(req entity.CreateVisitRequest) (*entity.VisitSingleResponse, []error)
	GetVisit(visitID string) (*entity.VisitSingleResponse, error)
	GetVisits(limit, offset int) (*entity.VisitMultiResponse, error)
	DeleteVisit(visitID string) error
}

type visitAdapter struct {
	visitUseCase usecase.VisitUseCase
	userAdapter  UserAdapter
}

func NewVisitAdapter(visitUseCase usecase.VisitUseCase, userAdapter UserAdapter) VisitAdapter {
	return &visitAdapter{
		visitUseCase: visitUseCase,
		userAdapter:  userAdapter,
	}
}

func (v *visitAdapter) CreateVisit(req entity.CreateVisitRequest) (*entity.VisitSingleResponse, []error) {
	visit, err := v.visitUseCase.CreateVisit(req)
	if err != nil {
		return nil, err
	}

	return v.setSingleVisitResponse(visit), nil
}

func (v *visitAdapter) GetVisit(visitID string) (*entity.VisitSingleResponse, error) {
	visitUUID, err := uuid.FromString(visitID)
	if err != nil {
		return nil, err
	}

	visit, err := v.visitUseCase.GetVisit(visitUUID)
	if err != nil {
		return nil, err
	}

	return v.setSingleVisitResponse(visit), nil
}

func (v *visitAdapter) GetVisits(limit, offset int) (*entity.VisitMultiResponse, error) {
	visits, err := v.visitUseCase.GetVisits(limit, offset)
	if err != nil {
		return nil, err
	}

	return v.setMultiVisitResponse(visits, limit, offset), nil
}

func (v *visitAdapter) DeleteVisit(visitID string) error {
	visitUUID, err := uuid.FromString(visitID)
	if err != nil {
		return err
	}

	return v.visitUseCase.DeleteVisit(visitUUID)
}

func (v *visitAdapter) setSingleVisitResponse(visit *entity.Visit) *entity.VisitSingleResponse {
	user, _ := v.userAdapter.GetUser(visit.UserID.String())

	resp := entity.VisitSingleResponse{
		ID:        visit.ID,
		User:      *user,
		CreatedAt: visit.CreatedAt,
		UpdatedAt: visit.UpdatedAt,
		DeletedAt: visit.DeletedAt,
	}

	return &resp
}

func (v *visitAdapter) setMultiVisitResponse(visits []entity.Visit, limit, offset int) *entity.VisitMultiResponse {
	var visitsSingleResp []entity.VisitSingleResponse
	var totalRows, totalPage, page int64
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, visit := range visits {
			visitsSingleResp = append(visitsSingleResp, *v.setSingleVisitResponse(&visit))
		}
	}()

	go func() {
		defer wg.Done()
		totalRows, _ = v.visitUseCase.CountVisit()
		totalPage, page = util.CountTotalPageAndCurrentPage(totalRows, limit, offset)
	}()

	wg.Wait()

	resp := entity.VisitMultiResponse{
		Page:      page,
		Offset:    offset,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
		Filter: entity.VisitQueryFilter{
			OrderBy: entity.VisitOrderBy{
				Field: "created_at",
				Sort:  "ascending",
			},
		},
		Rows: visitsSingleResp,
	}

	return &resp
}
