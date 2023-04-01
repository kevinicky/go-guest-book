package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/repository"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"time"
)

type VisitUseCase interface {
	CreateVisit(req entity.CreateVisitRequest) (*entity.Visit, []error)
	GetVisit(visitID uuid.UUID) (*entity.Visit, error)
	GetVisits(limit, offset int) ([]entity.Visit, error)
	CountVisit() (int64, error)
	DeleteVisit(visitID uuid.UUID) error
}

type visitUseCase struct {
	visitRepository repository.VisitRepository
	userUseCase     UserUseCase
}

func NewVisitUseCase(visitRepository repository.VisitRepository, userUseCase UserUseCase) VisitUseCase {
	return &visitUseCase{
		visitRepository: visitRepository,
		userUseCase:     userUseCase,
	}
}

func (v *visitUseCase) GetVisit(visitID uuid.UUID) (*entity.Visit, error) {
	keyCache := "visit~" + visitID.String()
	ctx := context.Background()
	visitCache := entity.Visit{}
	res, errCache := v.visitRepository.GetCacheData(ctx, keyCache)

	if errCache != nil {
		visit, err := v.visitRepository.FindVisit(visitID)
		if err != nil {
			return nil, err
		}

		resMarshall, _ := json.Marshal(visit)
		_, _ = v.visitRepository.SetCacheData(ctx, keyCache, resMarshall)

		return visit, err
	}

	_ = json.Unmarshal([]byte(res), &visitCache)

	return &visitCache, nil
}

func (v *visitUseCase) DeleteVisit(visitID uuid.UUID) error {
	ctx := context.Background()
	keyCache := "visit~" + visitID.String()
	_, _ = v.visitRepository.DeleteCacheData(ctx, keyCache)

	_, err := v.visitRepository.FindVisit(visitID)
	if err != nil {
		return err
	}

	return v.visitRepository.SoftDeleteVisit(visitID)
}

func (v *visitUseCase) GetVisits(limit, offset int) ([]entity.Visit, error) {
	return v.visitRepository.GetVisits(limit, offset)
}

func (v *visitUseCase) CountVisit() (int64, error) {
	return v.visitRepository.CountVisit()
}

func (v *visitUseCase) CreateVisit(req entity.CreateVisitRequest) (*entity.Visit, []error) {
	if req.UserID == uuid.Nil {
		return nil, []error{errors.New(customerror.USER_ID_IS_MANDATORY)}
	}

	_, err := v.userUseCase.GetUser(req.UserID, "")
	if err != nil {
		return nil, []error{err}
	}

	newID, _ := uuid.NewV4()
	visit := entity.Visit{
		ID:        newID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
	}

	visitResp, err := v.visitRepository.CreateVisit(visit)
	if err != nil {
		return nil, []error{err}
	}

	ctx := context.Background()
	keyCache := "visit~" + visitResp.ID.String()
	resMarshall, _ := json.Marshal(visitResp)
	_, _ = v.visitRepository.SetCacheData(ctx, keyCache, resMarshall)

	return visitResp, nil
}
