package product

import (
	"stock-service/internal/domain/product"
)

type CreateProductInput struct {
	TitleFa          string
	TitleEn          *string
	Slug             string
	Description      *string
	BrandID          int64
	CategoryID       int64
	OwnerType        product.OwnerType
	OwnerID          *int64
	IsOriginal       *bool
	MetaTitle        *string
	MetaDescription  *string
	IndexImageFileID *int64
}

type CreateProductUseCase struct {
	repo product.Repository
}

func NewCreateProductUseCase(repo product.Repository) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo}
}

func (uc *CreateProductUseCase) Execute(input CreateProductInput) (*product.Product, error) {
	var opts []product.ProductOption

	if input.TitleEn != nil {
		opts = append(opts, product.WithTitleEn(input.TitleEn))
	}
	if input.Description != nil {
		opts = append(opts, product.WithDescription(input.Description))
	}
	if input.OwnerType != "" {
		opts = append(opts, product.WithOwnerType(input.OwnerType))
	}
	if input.OwnerID != nil {
		opts = append(opts, product.WithOwnerID(input.OwnerID))
	}
	if input.IsOriginal != nil {
		opts = append(opts, product.WithIsOriginal(*input.IsOriginal))
	}
	if input.IndexImageFileID != nil {
		opts = append(opts, product.WithIndexImageFileID(input.IndexImageFileID))
	}
	if input.Slug != "" {
		opts = append(opts, product.WithSlug(input.Slug))
	}
	if input.MetaTitle != nil {
		opts = append(opts, product.WithMetaTitle(input.MetaTitle))
	}
	if input.MetaDescription != nil {
		opts = append(opts, product.WithMetaDescription(input.MetaDescription))
	}

	p, err := product.NewProduct(input.TitleFa, input.BrandID, input.CategoryID, opts...)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}

	return p, nil
}
