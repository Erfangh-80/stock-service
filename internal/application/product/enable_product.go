package product

import (
	"stock-service/internal/domain/product"
)

type EnableProductInput struct {
	ID int32
}

type EnableProductUseCase struct {
	repo product.Repository
}

func NewEnableProductUseCase(repo product.Repository) *EnableProductUseCase {
	return &EnableProductUseCase{repo: repo}
}

func (uc *EnableProductUseCase) Execute(input EnableProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.MarkEnabled()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
