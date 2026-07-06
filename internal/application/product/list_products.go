package product

import (
	"stock-service/internal/domain/product"
)

type ListProductsInput struct {
	OwnerType  *product.OwnerType
	OwnerID    *int64
	Status     *product.ProductStatus
	CategoryID *int64
	BrandID    *int64
	Search     *string
	Page       int
	Limit      int
}

type ListProductsOutput struct {
	Products []*product.Product
	Total    int
	Page     int
	Limit    int
}

type ListProductsUseCase struct {
	repo product.Repository
}

func NewListProductsUseCase(repo product.Repository) *ListProductsUseCase {
	return &ListProductsUseCase{repo: repo}
}

func (uc *ListProductsUseCase) Execute(input ListProductsInput) (*ListProductsOutput, error) {
	filter := product.ProductFilter{
		OwnerType:  input.OwnerType,
		OwnerID:    input.OwnerID,
		Status:     input.Status,
		CategoryID: input.CategoryID,
		BrandID:    input.BrandID,
		Search:     input.Search,
		Page:       input.Page,
		Limit:      input.Limit,
	}

	total, err := uc.repo.Count(filter)
	if err != nil {
		return nil, err
	}

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	products, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return &ListProductsOutput{
		Products: products,
		Total:    total,
		Page:     input.Page,
		Limit:    input.Limit,
	}, nil
}
