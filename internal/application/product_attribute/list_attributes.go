package productattribute

import (
	productattributedomain "stock-service/internal/domain/product_attribute"
)

type ListAttributesInput struct {
	ProductID int32
}

type ListAttributesOutput struct {
	Attributes []*productattributedomain.ProductAttribute
}

type ListAttributesUseCase struct {
	repo productattributedomain.Repository
}

func NewListAttributesUseCase(repo productattributedomain.Repository) *ListAttributesUseCase {
	return &ListAttributesUseCase{repo: repo}
}

func (uc *ListAttributesUseCase) Execute(input ListAttributesInput) (*ListAttributesOutput, error) {
	attrs, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	return &ListAttributesOutput{Attributes: attrs}, nil
}
