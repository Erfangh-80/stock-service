package warehouse

import (
	"stock-service/internal/domain/warehouse"
)

type ListWarehousesUseCase struct {
	repo warehouse.Repository
}

func NewListWarehousesUseCase(repo warehouse.Repository) *ListWarehousesUseCase {
	return &ListWarehousesUseCase{repo: repo}
}

type ListWarehousesInput struct {
	CreatedByUserID *int64
	IsPublic        *bool
	Page            int
	Limit           int
}

type ListWarehousesOutput struct {
	Warehouses []*warehouse.Warehouse
	Total      int
}

func (uc *ListWarehousesUseCase) Execute(input ListWarehousesInput) (*ListWarehousesOutput, error) {
	filter := warehouse.WarehouseFilter{
		CreatedByUserID: input.CreatedByUserID,
		IsPublic:        input.IsPublic,
		Page:            input.Page,
		Limit:           input.Limit,
	}
	items, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListWarehousesOutput{Warehouses: items, Total: total}, nil
}
