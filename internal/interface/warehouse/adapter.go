package warehouse

import (
	"errors"
	"stock-service/internal/domain/warehouse"
	iface "stock-service/internal/interface"
)



type CreateWarehouseInput struct {
	CreatedByUserID int64
	WarehouseName   string
}

type UpdateVisibilityInput struct {
	WarehouseID int64
	IsPublic    bool
}

type UpdateContactInput struct {
	WarehouseID      int64
	Phone            *string
	ContactPhone     *string
	CollectionMethod string
}

type WarehouseResponse struct {
	ID               int64  `json:"id"`
	CreatedByUserID  int64  `json:"created_by_user_id"`
	WarehouseName    string `json:"warehouse_name"`
	IsPublic         bool   `json:"is_public"`
	CollectionMethod string `json:"collection_method"`
}

type CreateWarehouseUseCase interface {
	Execute(input CreateWarehouseInput) (*warehouse.Warehouse, error)
}

type UpdateVisibilityUseCase interface {
	Execute(input UpdateVisibilityInput) (*warehouse.Warehouse, error)
}

type UpdateContactUseCase interface {
	Execute(input UpdateContactInput) (*warehouse.Warehouse, error)
}

type Adapter struct {
	create     CreateWarehouseUseCase
	updateVis  UpdateVisibilityUseCase
	updateCont UpdateContactUseCase
}

func NewAdapter(create CreateWarehouseUseCase, updateVis UpdateVisibilityUseCase, updateCont UpdateContactUseCase) *Adapter {
	return &Adapter{create: create, updateVis: updateVis, updateCont: updateCont}
}

func (a *Adapter) Create(input CreateWarehouseInput) (*WarehouseResponse, error) {
	w, err := a.create.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, warehouse.ErrWarehouseNameRequired),
			errors.Is(err, warehouse.ErrWarehouseNameTooLong):
			return nil, iface.ErrInvalidInput
		default:
			return nil, iface.ErrInternal
		}
	}
	return toResponse(w), nil
}

func (a *Adapter) UpdateVisibility(input UpdateVisibilityInput) (*WarehouseResponse, error) {
	w, err := a.updateVis.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, warehouse.ErrWarehouseNameRequired),
			errors.Is(err, warehouse.ErrWarehouseNameTooLong):
			return nil, iface.ErrInvalidInput
		default:
			return nil, iface.ErrInternal
		}
	}
	return toResponse(w), nil
}

func (a *Adapter) UpdateContact(input UpdateContactInput) (*WarehouseResponse, error) {
	w, err := a.updateCont.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, warehouse.ErrWarehouseNameRequired),
			errors.Is(err, warehouse.ErrWarehouseNameTooLong):
			return nil, iface.ErrInvalidInput
		default:
			return nil, iface.ErrInternal
		}
	}
	return toResponse(w), nil
}

func toResponse(w *warehouse.Warehouse) *WarehouseResponse {
	return &WarehouseResponse{
		ID:               w.ID,
		CreatedByUserID:  w.CreatedByUserID,
		WarehouseName:    w.WarehouseName,
		IsPublic:         w.IsPublic,
		CollectionMethod: w.CollectionMethod,
	}
}
