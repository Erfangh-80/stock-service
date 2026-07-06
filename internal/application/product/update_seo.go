package product

import (
	"stock-service/internal/domain/product"
)

type UpdateSEOInput struct {
	ID             int32
	MetaTitle      *string
	MetaDescription *string
}

type UpdateSEOUseCase struct {
	repo product.Repository
}

func NewUpdateSEOUseCase(repo product.Repository) *UpdateSEOUseCase {
	return &UpdateSEOUseCase{repo: repo}
}

func (uc *UpdateSEOUseCase) Execute(input UpdateSEOInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	p.UpdateSEO(input.MetaTitle, input.MetaDescription)

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
