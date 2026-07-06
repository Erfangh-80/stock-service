package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type DeleteWarehouseUseCase struct {
	repo warehouse.Repository
}

func NewDeleteWarehouseUseCase(repo warehouse.Repository) *DeleteWarehouseUseCase {
	return &DeleteWarehouseUseCase{repo: repo}
}

func (uc *DeleteWarehouseUseCase) Execute(warehouseID int64) error {
	w, err := uc.repo.FindByID(warehouseID)
	if err != nil {
		return err
	}
	if w == nil {
		return warehouse.ErrWarehouseNotFound
	}
	return uc.repo.Delete(warehouseID)
}
