package product

import (
	"stock-service/internal/domain/product"
)

type SoftDeleteProductInput struct {
	ID int32
}

type SoftDeleteProductUseCase struct {
	repo product.Repository
}

func NewSoftDeleteProductUseCase(repo product.Repository) *SoftDeleteProductUseCase {
	return &SoftDeleteProductUseCase{repo: repo}
}

func (uc *SoftDeleteProductUseCase) Execute(input SoftDeleteProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.SoftDelete()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
