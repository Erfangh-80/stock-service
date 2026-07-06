package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type UpdateWarehouseUseCase struct {
	repo warehouse.Repository
}

func NewUpdateWarehouseUseCase(repo warehouse.Repository) *UpdateWarehouseUseCase {
	return &UpdateWarehouseUseCase{repo: repo}
}

type UpdateWarehouseInput struct {
	WarehouseID int64
	Name        *string
	AddressID   *int64
}

func (uc *UpdateWarehouseUseCase) Execute(input UpdateWarehouseInput) (*warehouse.Warehouse, error) {
	w, err := uc.repo.FindByID(input.WarehouseID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, warehouse.ErrWarehouseNotFound
	}
	if input.Name != nil {
		if err := w.UpdateWarehouseName(*input.Name); err != nil {
			return nil, err
		}
	}
	if input.AddressID != nil {
		if err := w.UpdateAddressID(*input.AddressID); err != nil {
			return nil, err
		}
	}
	if err := uc.repo.Save(w); err != nil {
		return nil, err
	}
	return w, nil
}
