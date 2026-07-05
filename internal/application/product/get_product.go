package product

import (
	"stock-service/internal/domain/product"
)

type GetProductInput struct {
	ID int32
}

type GetProductUseCase struct {
	repo product.Repository
}

func NewGetProductUseCase(repo product.Repository) *GetProductUseCase {
	return &GetProductUseCase{repo: repo}
}

func (uc *GetProductUseCase) Execute(input GetProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}
	return p, nil
}
