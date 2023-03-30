package usecase

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"strings"
	"time"
)

type ThreadUseCase interface {
	CreateThread(req entity.CreateThreadRequest) (*entity.Thread, []error)
	GetThread(threadID uuid.UUID) (*entity.Thread, error)
	GetThreads(limit, offset int, visitID uuid.UUID) ([]entity.Thread, error)
	CountThread(visitID uuid.UUID) (int64, error)
	DeleteThread(threadID uuid.UUID) error
}

type threadUseCase struct {
	threadRepository repository.ThreadRepository
	visitUseCase     VisitUseCase
	userUseCase      UserUseCase
}

func NewThreadUseCase(threadRepository repository.ThreadRepository, visitUseCase VisitUseCase, userUseCase UserUseCase) ThreadUseCase {
	return &threadUseCase{
		threadRepository: threadRepository,
		visitUseCase:     visitUseCase,
		userUseCase:      userUseCase,
	}
}

func (t *threadUseCase) GetThread(threadID uuid.UUID) (*entity.Thread, error) {
	return t.threadRepository.FindThread(threadID)
}

func (t *threadUseCase) DeleteThread(threadID uuid.UUID) error {
	_, err := t.threadRepository.FindThread(threadID)
	if err != nil {
		return err
	}

	return t.threadRepository.SoftDeleteThread(threadID)
}

func (t *threadUseCase) GetThreads(limit, offset int, visitID uuid.UUID) ([]entity.Thread, error) {
	return t.threadRepository.GetThreads(limit, offset, visitID)
}

func (t *threadUseCase) CountThread(visitID uuid.UUID) (int64, error) {
	return t.threadRepository.CountThread(visitID)
}

func (t *threadUseCase) CreateThread(req entity.CreateThreadRequest) (*entity.Thread, []error) {
	var errList []error
	sanitiseContent := strings.TrimSpace(req.Content)

	if req.VisitID == uuid.Nil {
		errList = append(errList, errors.New(customerror.VISIT_ID_IS_MANDATORY))
	} else {
		_, err := t.visitUseCase.GetVisit(req.VisitID)
		if err != nil {
			errList = append(errList, err)
		}
	}

	if req.UserID == uuid.Nil {
		errList = append(errList, errors.New(customerror.USER_ID_IS_MANDATORY))
	} else {
		_, err := t.userUseCase.GetUser(req.UserID)
		if err != nil {
			errList = append(errList, err)
		}
	}

	if sanitiseContent == "" {
		errList = append(errList, errors.New(customerror.CONTENT_IS_MANDATORY))
	}

	if len(errList) > 0 {
		return nil, errList
	}

	newID, _ := uuid.NewV4()
	thread := entity.Thread{
		ID:        newID,
		VisitID:   req.VisitID,
		UserID:    req.UserID,
		Content:   sanitiseContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}

	threadResp, err := t.threadRepository.CreateThread(thread)
	if err != nil {
		return nil, []error{err}
	}
	return threadResp, nil
}
