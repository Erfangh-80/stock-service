package product

import (
	"stock-service/internal/domain/product"
)

type ActivateProductInput struct {
	ID int32
}

type ActivateProductUseCase struct {
	repo product.Repository
}

func NewActivateProductUseCase(repo product.Repository) *ActivateProductUseCase {
	return &ActivateProductUseCase{repo: repo}
}

func (uc *ActivateProductUseCase) Execute(input ActivateProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.MarkActive()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
