package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type GetWarehouseUseCase struct {
	repo warehouse.Repository
}

func NewGetWarehouseUseCase(repo warehouse.Repository) *GetWarehouseUseCase {
	return &GetWarehouseUseCase{repo: repo}
}

func (uc *GetWarehouseUseCase) Execute(warehouseID int64) (*warehouse.Warehouse, error) {
	w, err := uc.repo.FindByID(warehouseID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, warehouse.ErrWarehouseNotFound
	}
	return w, nil
}
