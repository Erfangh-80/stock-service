package storeallowedcategoryinterface

import (
	"time"

	domain "stock-service/internal/domain/store_allowed_category"
	iface "stock-service/internal/interface"
)

type StoreAllowedCategoryOutput struct {
	ID          int64
	StoreID     int64
	CategoryID  int64
	Status      string
	SupportNote string
	CreatedAt   time.Time
}

type CreateCategoryParams struct {
	StoreID    int64
	CategoryID int64
}

type ApproveCategoryParams struct {
	CategoryID int64
}

type RejectCategoryParams struct {
	CategoryID int64
}

type createCategoryUseCase interface {
	Execute(storeID, categoryID int64) (*domain.StoreAllowedCategory, error)
}

type approveCategoryUseCase interface {
	Execute(categoryID int64) error
}

type rejectCategoryUseCase interface {
	Execute(categoryID int64) error
}

type Adapter struct {
	create  createCategoryUseCase
	approve approveCategoryUseCase
	reject  rejectCategoryUseCase
}

func NewAdapter(
	create createCategoryUseCase,
	approve approveCategoryUseCase,
	reject rejectCategoryUseCase,
) *Adapter {
	return &Adapter{
		create:  create,
		approve: approve,
		reject:  reject,
	}
}

func (a *Adapter) Create(params CreateCategoryParams) (*StoreAllowedCategoryOutput, error) {
	result, err := a.create.Execute(params.StoreID, params.CategoryID)
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) Approve(params ApproveCategoryParams) error {
	err := a.approve.Execute(params.CategoryID)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Reject(params RejectCategoryParams) error {
	err := a.reject.Execute(params.CategoryID)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	return iface.ErrInternal
}

func toOutput(sac *domain.StoreAllowedCategory) *StoreAllowedCategoryOutput {
	return &StoreAllowedCategoryOutput{
		ID:          sac.ID,
		StoreID:     sac.StoreID,
		CategoryID:  sac.CategoryID,
		Status:      string(sac.Status),
		SupportNote: sac.SupportNote,
		CreatedAt:   sac.CreatedAt,
	}
}
