package storewarehouselink

import (
	"errors"

	app "stock-service/internal/application/store_warehouse_link"
	domain "stock-service/internal/domain/store_warehouse_link"
	iface "stock-service/internal/interface"
)

type CreateLinkInput struct {
	StoreID     int64 `json:"store_id"`
	WarehouseID int64 `json:"warehouse_id"`
}

type GetLinkInput struct {
	ID int64
}

type ListLinksInput struct {
	StoreID     *int64
	WarehouseID *int64
	Page        int
	Limit       int
}

type DeleteLinkInput struct {
	ID int64
}

type ChangeRelationInput struct {
	LinkID       int64
	RelationType string
}

type LinkResponse struct {
	ID           int64  `json:"id"`
	StoreID      int64  `json:"store_id"`
	WarehouseID  int64  `json:"warehouse_id"`
	RelationType string `json:"relation_type"`
}

type LinkListResponse struct {
	Links []LinkResponse `json:"links"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type createLinkUseCase interface {
	Execute(storeID, warehouseID int64) (*domain.StoreWarehouseLink, error)
}

type getLinkUseCase interface {
	Execute(input app.GetLinkInput) (*domain.StoreWarehouseLink, error)
}

type listLinksUseCase interface {
	Execute(input app.ListLinksInput) (*app.ListLinksOutput, error)
}

type changeRelationUseCase interface {
	Execute(input app.ChangeRelationInput) (*domain.StoreWarehouseLink, error)
}

type deleteLinkUseCase interface {
	Execute(input app.DeleteLinkInput) error
}

type Adapter struct {
	create createLinkUseCase
	get    getLinkUseCase
	list   listLinksUseCase
	change changeRelationUseCase
	delete deleteLinkUseCase
}

func NewAdapter(
	create createLinkUseCase,
	get getLinkUseCase,
	list listLinksUseCase,
	change changeRelationUseCase,
	del deleteLinkUseCase,
) *Adapter {
	return &Adapter{
		create: create,
		get:    get,
		list:   list,
		change: change,
		delete: del,
	}
}

func (a *Adapter) Create(input CreateLinkInput) (*LinkResponse, error) {
	swl, err := a.create.Execute(input.StoreID, input.WarehouseID)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(swl), nil
}

func (a *Adapter) Get(input GetLinkInput) (*LinkResponse, error) {
	swl, err := a.get.Execute(app.GetLinkInput{ID: input.ID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(swl), nil
}

func (a *Adapter) List(input ListLinksInput) (*LinkListResponse, error) {
	result, err := a.list.Execute(app.ListLinksInput{
		StoreID:     input.StoreID,
		WarehouseID: input.WarehouseID,
		Page:        input.Page,
		Limit:       input.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]LinkResponse, len(result.Links))
	for i, swl := range result.Links {
		items[i] = *toResponse(swl)
	}

	return &LinkListResponse{
		Links: items,
		Total: result.Total,
		Page:  result.Page,
		Limit: result.Limit,
	}, nil
}

func (a *Adapter) ChangeRelation(input ChangeRelationInput) (*LinkResponse, error) {
	swl, err := a.change.Execute(app.ChangeRelationInput{
		LinkID:       input.LinkID,
		RelationType: domain.RelationType(input.RelationType),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(swl), nil
}

func (a *Adapter) Delete(input DeleteLinkInput) error {
	err := a.delete.Execute(app.DeleteLinkInput{ID: input.ID})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrLinkNotFound):
		return iface.ErrNotFound
	case errors.Is(err, domain.ErrInvalidRelationType):
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(swl *domain.StoreWarehouseLink) *LinkResponse {
	return &LinkResponse{
		ID:           swl.ID,
		StoreID:      swl.StoreID,
		WarehouseID:  swl.WarehouseID,
		RelationType: string(swl.RelationType),
	}
}
