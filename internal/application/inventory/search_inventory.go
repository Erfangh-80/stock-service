package inventory

import (
	"stock-service/internal/domain/inventory"
	productdomain "stock-service/internal/domain/product"
)

type SearchInventoryInput struct {
	Query string
	Page  int
	Limit int
}

type SearchInventoryOutput struct {
	Items []*inventory.Inventory
	Total int
	Page  int
	Limit int
}

type SearchInventoryUseCase interface {
	Execute(input SearchInventoryInput) (*SearchInventoryOutput, error)
}

type searchInventoryInteractor struct {
	invRepo     inventory.Repository
	productRepo productdomain.Repository
}

func NewSearchInventoryUseCase(invRepo inventory.Repository, productRepo productdomain.Repository) SearchInventoryUseCase {
	return &searchInventoryInteractor{invRepo: invRepo, productRepo: productRepo}
}

func (uc *searchInventoryInteractor) Execute(input SearchInventoryInput) (*SearchInventoryOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 || input.Limit > 100 {
		input.Limit = 20
	}

	products, err := uc.productRepo.FindByTitle(input.Query)
	if err != nil {
		return nil, err
	}

	productIDs := make(map[int32]struct{}, len(products))
	for _, p := range products {
		productIDs[p.ID] = struct{}{}
	}

	all, err := uc.invRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var filtered []*inventory.Inventory
	for _, inv := range all {
		if _, ok := productIDs[inv.ProductID]; ok {
			filtered = append(filtered, inv)
		}
	}

	total := len(filtered)
	start := (input.Page - 1) * input.Limit
	if start >= total {
		return &SearchInventoryOutput{Items: []*inventory.Inventory{}, Total: total, Page: input.Page, Limit: input.Limit}, nil
	}

	end := start + input.Limit
	if end > total {
		end = total
	}

	return &SearchInventoryOutput{Items: filtered[start:end], Total: total, Page: input.Page, Limit: input.Limit}, nil
}
