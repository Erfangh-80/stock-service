package product

import (
	"stock-service/internal/domain/product"
)

type UpdateProductInput struct {
	ID               int32
	TitleFa          *string
	TitleEn          *string
	Description      *string
	BrandID          *int64
	CategoryID       *int64
	IndexImageFileID *int64
}

type UpdateProductUseCase struct {
	repo product.Repository
}

func NewUpdateProductUseCase(repo product.Repository) *UpdateProductUseCase {
	return &UpdateProductUseCase{repo: repo}
}

func (uc *UpdateProductUseCase) Execute(input UpdateProductInput) (*product.Product, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	if input.TitleFa != nil {
		if *input.TitleFa == "" {
			return nil, product.ErrTitleFaRequired
		}
		p.TitleFa = *input.TitleFa
	}
	if input.TitleEn != nil {
		p.TitleEn = input.TitleEn
	}
	if input.Description != nil {
		p.Description = input.Description
	}
	if input.BrandID != nil {
		if *input.BrandID <= 0 {
			return nil, product.ErrInvalidBrandID
		}
		p.BrandID = *input.BrandID
	}
	if input.CategoryID != nil {
		if *input.CategoryID <= 0 {
			return nil, product.ErrInvalidCategoryID
		}
		p.CategoryID = *input.CategoryID
	}
	if input.IndexImageFileID != nil {
		p.IndexImageFileID = input.IndexImageFileID
	}

	p.Touch()

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
