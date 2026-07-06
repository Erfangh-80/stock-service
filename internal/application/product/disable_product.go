package product

import (
	"stock-service/internal/domain/product"
)

type DisableProductInput struct {
	ID int32
}

type DisableProductUseCase struct {
	repo product.Repository
}

func NewDisableProductUseCase(repo product.Repository) *DisableProductUseCase {
	return &DisableProductUseCase{repo: repo}
}

func (uc *DisableProductUseCase) Execute(input DisableProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.MarkDisabled()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
