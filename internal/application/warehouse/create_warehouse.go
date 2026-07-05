package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type CreateWarehouseUseCase struct {
	repo warehouse.Repository
}

func NewCreateWarehouseUseCase(repo warehouse.Repository) *CreateWarehouseUseCase {
	return &CreateWarehouseUseCase{repo: repo}
}

func (uc *CreateWarehouseUseCase) Execute(createdByUserID int64, warehouseName string) (*warehouse.Warehouse, error) {
	w, err := warehouse.NewWarehouse(createdByUserID, warehouseName)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(w); err != nil {
		return nil, err
	}
	return w, nil
}
