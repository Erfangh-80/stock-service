package warehouse

import (
	"errors"

	app "stock-service/internal/application/warehouse"
	domain "stock-service/internal/domain/warehouse"
	iface "stock-service/internal/interface"
)

type CreateWarehouseInput struct {
	CreatedByUserID int64
	WarehouseName   string
}

type GetWarehouseInput struct {
	ID int64
}

type ListWarehousesInput struct {
	CreatedByUserID *int64
	IsPublic        *bool
	Page            int
	Limit           int
}

type DeleteWarehouseInput struct {
	ID int64
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

type UpdateWarehouseInput struct {
	WarehouseID int64
	Name        *string
	AddressID   *int64
}

type WarehouseResponse struct {
	ID               int64   `json:"id"`
	CreatedByUserID  int64   `json:"created_by_user_id"`
	WarehouseName    string  `json:"warehouse_name"`
	AddressID        *int64  `json:"address_id,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	ContactPhone     *string `json:"contact_phone,omitempty"`
	IsPublic         bool    `json:"is_public"`
	CollectionMethod string  `json:"collection_method"`
	CreatedAt        string  `json:"created_at"`
}

type WarehouseListResponse struct {
	Warehouses []WarehouseResponse `json:"warehouses"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}

type createWarehouseUseCase interface {
	Execute(createdByUserID int64, warehouseName string) (*domain.Warehouse, error)
}

type getWarehouseUseCase interface {
	Execute(warehouseID int64) (*domain.Warehouse, error)
}

type listWarehousesUseCase interface {
	Execute(input app.ListWarehousesInput) (*app.ListWarehousesOutput, error)
}

type deleteWarehouseUseCase interface {
	Execute(warehouseID int64) error
}

type updateVisibilityUseCase interface {
	Execute(warehouseID int64, isPublic bool) error
}

type updateContactUseCase interface {
	Execute(warehouseID int64, phone, contactPhone *string, collectionMethod string) error
}

type updateWarehouseUseCase interface {
	Execute(input app.UpdateWarehouseInput) (*domain.Warehouse, error)
}

type Adapter struct {
	create     createWarehouseUseCase
	get        getWarehouseUseCase
	list       listWarehousesUseCase
	del        deleteWarehouseUseCase
	updateVis  updateVisibilityUseCase
	updateCont updateContactUseCase
	updateWH   updateWarehouseUseCase
}

func NewAdapter(
	create createWarehouseUseCase,
	get getWarehouseUseCase,
	list listWarehousesUseCase,
	del deleteWarehouseUseCase,
	updateVis updateVisibilityUseCase,
	updateCont updateContactUseCase,
	updateWH updateWarehouseUseCase,
) *Adapter {
	return &Adapter{
		create:     create,
		get:        get,
		list:       list,
		del:        del,
		updateVis:  updateVis,
		updateCont: updateCont,
		updateWH:   updateWH,
	}
}

func (a *Adapter) Create(input CreateWarehouseInput) (*WarehouseResponse, error) {
	w, err := a.create.Execute(input.CreatedByUserID, input.WarehouseName)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(w), nil
}

func (a *Adapter) Get(input GetWarehouseInput) (*WarehouseResponse, error) {
	w, err := a.get.Execute(input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(w), nil
}

func (a *Adapter) List(input ListWarehousesInput) (*WarehouseListResponse, error) {
	result, err := a.list.Execute(app.ListWarehousesInput{
		CreatedByUserID: input.CreatedByUserID,
		IsPublic:        input.IsPublic,
		Page:            input.Page,
		Limit:           input.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}

	items := make([]WarehouseResponse, len(result.Warehouses))
	for i, w := range result.Warehouses {
		items[i] = *toResponse(w)
	}

	page := input.Page
	if page < 1 {
		page = 1
	}
	limit := input.Limit
	if limit < 1 {
		limit = 20
	}

	return &WarehouseListResponse{
		Warehouses: items,
		Total:      result.Total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (a *Adapter) Delete(input DeleteWarehouseInput) error {
	err := a.del.Execute(input.ID)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) UpdateVisibility(input UpdateVisibilityInput) (*WarehouseResponse, error) {
	if err := a.updateVis.Execute(input.WarehouseID, input.IsPublic); err != nil {
		return nil, mapError(err)
	}
	w, err := a.get.Execute(input.WarehouseID)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(w), nil
}

func (a *Adapter) UpdateContact(input UpdateContactInput) (*WarehouseResponse, error) {
	if err := a.updateCont.Execute(input.WarehouseID, input.Phone, input.ContactPhone, input.CollectionMethod); err != nil {
		return nil, mapError(err)
	}
	w, err := a.get.Execute(input.WarehouseID)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(w), nil
}

func (a *Adapter) UpdateWarehouse(input UpdateWarehouseInput) (*WarehouseResponse, error) {
	w, err := a.updateWH.Execute(app.UpdateWarehouseInput{
		WarehouseID: input.WarehouseID,
		Name:        input.Name,
		AddressID:   input.AddressID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(w), nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrWarehouseNotFound):
		return iface.ErrNotFound
	case errors.Is(err, domain.ErrWarehouseNameRequired),
		errors.Is(err, domain.ErrWarehouseNameTooLong),
		errors.Is(err, domain.ErrInvalidCollectionMethod),
		errors.Is(err, domain.ErrWarehouseAddressIDNotPositive):
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(w *domain.Warehouse) *WarehouseResponse {
	createdAt := ""
	if !w.CreatedAt.IsZero() {
		createdAt = w.CreatedAt.Format("2006-01-02T15:04:05Z")
	}
	return &WarehouseResponse{
		ID:               w.ID,
		CreatedByUserID:  w.CreatedByUserID,
		WarehouseName:    w.WarehouseName,
		AddressID:        w.AddressID,
		Phone:            w.Phone,
		ContactPhone:     w.ContactPhone,
		IsPublic:         w.IsPublic,
		CollectionMethod: w.CollectionMethod,
		CreatedAt:        createdAt,
	}
}
