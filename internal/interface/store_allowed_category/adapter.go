package storeallowedcategoryinterface

import (
	"errors"
	"time"

	app "stock-service/internal/application/store_allowed_category"
	domain "stock-service/internal/domain/store_allowed_category"
	iface "stock-service/internal/interface"
)

type StoreAllowedCategoryOutput struct {
	ID          int64     `json:"id"`
	StoreID     int64     `json:"store_id"`
	CategoryID  int64     `json:"category_id"`
	Status      string    `json:"status"`
	SupportNote string    `json:"support_note"`
	CreatedAt   time.Time `json:"created_at"`
}

type StoreCategoryListItem struct {
	ID          int64     `json:"id"`
	StoreID     int64     `json:"store_id"`
	CategoryID  int64     `json:"category_id"`
	Status      string    `json:"status"`
	SupportNote string    `json:"support_note"`
	CreatedAt   time.Time `json:"created_at"`
}

type StoreCategoryListResponse struct {
	Categories []StoreCategoryListItem `json:"categories"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
}

type CreateCategoryParams struct {
	StoreID    int64 `json:"store_id"`
	CategoryID int64 `json:"category_id"`
}

type GetCategoryParams struct {
	ID int64
}

type ListCategoriesParams struct {
	StoreID *int64
	Page    int
	Limit   int
}

type DeleteCategoryParams struct {
	ID int64
}

type RejectCategoryParams struct {
	CategoryID  int64
	SupportNote string
}

type createCategoryUseCase interface {
	Execute(storeID, categoryID int64) (*domain.StoreAllowedCategory, error)
}

type getCategoryUseCase interface {
	Execute(input app.GetStoreCategoryInput) (*domain.StoreAllowedCategory, error)
}

type listCategoriesUseCase interface {
	Execute(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error)
}

type approveCategoryUseCase interface {
	Execute(input app.ApproveCategoryInput) error
}

type rejectCategoryUseCase interface {
	Execute(input app.RejectCategoryInput) error
}

type deleteCategoryUseCase interface {
	Execute(input app.DeleteStoreCategoryInput) error
}

type validateCategoryExistsUseCase interface {
	Execute(input app.ValidateCategoryExistsInput) error
}

type Adapter struct {
	create    createCategoryUseCase
	get       getCategoryUseCase
	list      listCategoriesUseCase
	approve   approveCategoryUseCase
	reject    rejectCategoryUseCase
	delete    deleteCategoryUseCase
	validate  validateCategoryExistsUseCase
}

func NewAdapter(
	create createCategoryUseCase,
	get getCategoryUseCase,
	list listCategoriesUseCase,
	approve approveCategoryUseCase,
	reject rejectCategoryUseCase,
	delete deleteCategoryUseCase,
	validate validateCategoryExistsUseCase,
) *Adapter {
	return &Adapter{
		create:   create,
		get:      get,
		list:     list,
		approve:  approve,
		reject:   reject,
		delete:   delete,
		validate: validate,
	}
}

func (a *Adapter) Create(params CreateCategoryParams) (*StoreAllowedCategoryOutput, error) {
	if err := a.validate.Execute(app.ValidateCategoryExistsInput{CategoryID: params.CategoryID}); err != nil {
		return nil, mapError(err)
	}

	result, err := a.create.Execute(params.StoreID, params.CategoryID)
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) Get(params GetCategoryParams) (*StoreAllowedCategoryOutput, error) {
	result, err := a.get.Execute(app.GetStoreCategoryInput{ID: params.ID})
	if err != nil {
		return nil, mapError(err)
	}
	return toOutput(result), nil
}

func (a *Adapter) List(params ListCategoriesParams) (*StoreCategoryListResponse, error) {
	result, err := a.list.Execute(app.ListStoreCategoriesInput{
		StoreID: params.StoreID,
		Page:    params.Page,
		Limit:   params.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]StoreCategoryListItem, len(result.Categories))
	for i, cat := range result.Categories {
		items[i] = toListItem(cat)
	}

	return &StoreCategoryListResponse{
		Categories: items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
	}, nil
}

func (a *Adapter) Approve(params struct{ CategoryID int64 }) error {
	err := a.approve.Execute(app.ApproveCategoryInput{CategoryID: params.CategoryID})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Reject(params RejectCategoryParams) error {
	err := a.reject.Execute(app.RejectCategoryInput{
		CategoryID:  params.CategoryID,
		SupportNote: params.SupportNote,
	})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Delete(params DeleteCategoryParams) error {
	err := a.delete.Execute(app.DeleteStoreCategoryInput{ID: params.ID})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrStoreCategoryNotFound),
		errors.Is(err, domain.ErrCategoryNotFound):
		return iface.ErrNotFound
	case errors.Is(err, domain.ErrSupportNoteTooLong):
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
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

func toListItem(sac *domain.StoreAllowedCategory) StoreCategoryListItem {
	return StoreCategoryListItem{
		ID:          sac.ID,
		StoreID:     sac.StoreID,
		CategoryID:  sac.CategoryID,
		Status:      string(sac.Status),
		SupportNote: sac.SupportNote,
		CreatedAt:   sac.CreatedAt,
	}
}
