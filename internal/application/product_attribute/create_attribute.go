package productattribute

import (
	productattributedomain "stock-service/internal/domain/product_attribute"
)

type CreateAttributeInput struct {
	ProductID int32
	Key       string
	Value     string
}

type CreateAttributeUseCase struct {
	repo productattributedomain.Repository
}

func NewCreateAttributeUseCase(repo productattributedomain.Repository) *CreateAttributeUseCase {
	return &CreateAttributeUseCase{repo: repo}
}

func (uc *CreateAttributeUseCase) Execute(input CreateAttributeInput) (*productattributedomain.ProductAttribute, error) {
	pa, err := productattributedomain.NewProductAttribute(input.ProductID, input.Key, input.Value)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(pa); err != nil {
		return nil, err
	}
	return pa, nil
}
