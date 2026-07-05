package inventory

import (
	"stock-service/internal/domain/inventory"
)

type ListInventoryInput struct {
	StoreID          *int64
	ProductID        *int32
	VendorSaleStatus *string
	SystemSaleStatus *string
	Page             int
	Limit            int
}

type ListInventoryOutput struct {
	Items []*inventory.Inventory
	Total int
	Page  int
	Limit int
}

type ListInventoryUseCase interface {
	Execute(input ListInventoryInput) (*ListInventoryOutput, error)
}

type listInventoryInteractor struct {
	repo inventory.Repository
}

func NewListInventoryUseCase(repo inventory.Repository) ListInventoryUseCase {
	return &listInventoryInteractor{repo: repo}
}

func (uc *listInventoryInteractor) Execute(input ListInventoryInput) (*ListInventoryOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 || input.Limit > 100 {
		input.Limit = 20
	}

	all, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []*inventory.Inventory
	for _, inv := range all {
		if input.StoreID != nil && inv.StoreID != *input.StoreID {
			continue
		}
		if input.ProductID != nil && inv.ProductID != *input.ProductID {
			continue
		}
		if input.VendorSaleStatus != nil && string(inv.VendorSaleStatus) != *input.VendorSaleStatus {
			continue
		}
		if input.SystemSaleStatus != nil && string(inv.SystemSaleStatus) != *input.SystemSaleStatus {
			continue
		}
		filtered = append(filtered, inv)
	}

	total := len(filtered)
	start := (input.Page - 1) * input.Limit
	if start >= total {
		return &ListInventoryOutput{Items: []*inventory.Inventory{}, Total: total, Page: input.Page, Limit: input.Limit}, nil
	}

	end := start + input.Limit
	if end > total {
		end = total
	}

	return &ListInventoryOutput{Items: filtered[start:end], Total: total, Page: input.Page, Limit: input.Limit}, nil
}
