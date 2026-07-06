package promotioninterface

import (
	apppromotion "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
	iface "stock-service/internal/interface"
)

type PromotionResponse struct {
	ID        int64
	Title     string
	Status    string
	CreatedAt string
}

type CreatePromotionParams struct {
	Title string
}

type createPromotionUseCase interface {
	Execute(title string) (*promotion.Promotion, error)
}

type getPromotionUseCase interface {
	Execute(input apppromotion.GetPromotionInput) (*promotion.Promotion, error)
}

type activatePromotionUseCase interface {
	Execute(id int64) error
}

type deactivatePromotionUseCase interface {
	Execute(id int64) error
}

type Adapter struct {
	create     createPromotionUseCase
	get        getPromotionUseCase
	activate   activatePromotionUseCase
	deactivate deactivatePromotionUseCase
}

func NewAdapter(
	create createPromotionUseCase,
	get getPromotionUseCase,
	activate activatePromotionUseCase,
	deactivate deactivatePromotionUseCase,
) *Adapter {
	return &Adapter{create: create, get: get, activate: activate, deactivate: deactivate}
}

func (a *Adapter) Create(params CreatePromotionParams) (*PromotionResponse, error) {
	result, err := a.create.Execute(params.Title)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*PromotionResponse, error) {
	result, err := a.get.Execute(apppromotion.GetPromotionInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Activate(id int64) error {
	err := a.activate.Execute(id)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Deactivate(id int64) error {
	err := a.deactivate.Execute(id)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch err {
	case promotion.ErrPromotionNotFound:
		return iface.ErrNotFound
	case promotion.ErrTitleRequired, promotion.ErrTitleTooLong:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(p *promotion.Promotion) *PromotionResponse {
	return &PromotionResponse{
		ID: p.ID, Title: p.Title,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
