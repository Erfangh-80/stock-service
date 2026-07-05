package product

import (
	"stock-service/internal/domain/product"
)

type RejectProductInput struct {
	ID int32
}

type RejectProductUseCase struct {
	repo product.Repository
}

func NewRejectProductUseCase(repo product.Repository) *RejectProductUseCase {
	return &RejectProductUseCase{repo: repo}
}

func (uc *RejectProductUseCase) Execute(input RejectProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.MarkRejected()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
