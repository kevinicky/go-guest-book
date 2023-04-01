package adapter

import (
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/usecase"
	"github.com/kevinicky/go-guest-book/util"
	"sync"
)

type ThreadAdapter interface {
	CreateThread(req entity.CreateThreadRequest) (*entity.ThreadSingleResponse, []error)
	GetThread(threadID string) (*entity.ThreadSingleResponse, error)
	GetThreads(limit, offset int, visitID string) (*entity.ThreadMultiResponse, error)
	DeleteThread(threadID string) error
}

type threadAdapter struct {
	threadUseCase usecase.ThreadUseCase
	visitAdapter  VisitAdapter
	userAdapter   UserAdapter
}

func NewThreadAdapter(threadUseCase usecase.ThreadUseCase, visitAdapter VisitAdapter, userAdapter UserAdapter) ThreadAdapter {
	return &threadAdapter{
		threadUseCase: threadUseCase,
		visitAdapter:  visitAdapter,
		userAdapter:   userAdapter,
	}
}

func (v *threadAdapter) CreateThread(req entity.CreateThreadRequest) (*entity.ThreadSingleResponse, []error) {
	thread, err := v.threadUseCase.CreateThread(req)
	if err != nil {
		return nil, err
	}

	return v.setSingleThreadResponse(thread), nil
}

func (v *threadAdapter) GetThread(threadID string) (*entity.ThreadSingleResponse, error) {
	threadUUID, err := uuid.FromString(threadID)
	if err != nil {
		return nil, err
	}

	thread, err := v.threadUseCase.GetThread(threadUUID)
	if err != nil {
		return nil, err
	}

	return v.setSingleThreadResponse(thread), nil
}

func (v *threadAdapter) GetThreads(limit, offset int, visitID string) (*entity.ThreadMultiResponse, error) {
	visitUUID, err := uuid.FromString(visitID)
	if err != nil {
		return nil, err
	}

	threads, err := v.threadUseCase.GetThreads(limit, offset, visitUUID)
	if err != nil {
		return nil, err
	}

	return v.setMultiThreadResponse(threads, limit, offset, visitUUID), nil
}

func (v *threadAdapter) DeleteThread(threadID string) error {
	threadUUID, err := uuid.FromString(threadID)
	if err != nil {
		return err
	}

	return v.threadUseCase.DeleteThread(threadUUID)
}

func (v *threadAdapter) setSingleThreadResponse(thread *entity.Thread) *entity.ThreadSingleResponse {
	wg := sync.WaitGroup{}
	var visit *entity.VisitSingleResponse
	var user *entity.UserSingleResponse

	wg.Add(2)
	go func() {
		defer wg.Done()
		visit, _ = v.visitAdapter.GetVisit(thread.VisitID.String())
	}()
	go func() {
		defer wg.Done()
		user, _ = v.userAdapter.GetUser(thread.UserID.String())
	}()
	wg.Wait()

	resp := entity.ThreadSingleResponse{
		ID:        thread.ID,
		Visit:     *visit,
		User:      *user,
		Content:   thread.Content,
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
		DeletedAt: thread.DeletedAt,
	}

	return &resp
}

func (v *threadAdapter) setMultiThreadResponse(threads []entity.Thread, limit, offset int, visitID uuid.UUID) *entity.ThreadMultiResponse {
	var threadMultiResponsePerRow []entity.ThreadMultiResponsePerRow
	var totalRows, totalPage, page int64
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, thread := range threads {
			tempRow := *v.setSingleThreadResponse(&thread)
			threadMultiResponsePerRow = append(threadMultiResponsePerRow, entity.ThreadMultiResponsePerRow{
				ID:        tempRow.ID,
				User:      tempRow.User,
				Content:   tempRow.Content,
				CreatedAt: tempRow.CreatedAt,
				UpdatedAt: tempRow.UpdatedAt,
				DeletedAt: tempRow.DeletedAt,
			})
		}
	}()

	go func() {
		defer wg.Done()
		totalRows, _ = v.threadUseCase.CountThread(visitID)
		totalPage, page = util.CountTotalPageAndCurrentPage(totalRows, limit, offset)
	}()

	wg.Wait()

	resp := entity.ThreadMultiResponse{
		Page:      page,
		Offset:    offset,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
		Filter: entity.ThreadQueryFilter{
			VisitID: visitID,
			OrderBy: entity.ThreadOrderBy{
				Field: "created_at",
				Sort:  "ascending",
			},
		},
		Rows: threadMultiResponsePerRow,
	}

	return &resp
}
