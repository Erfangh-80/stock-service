package producttype

import (
	producttypedomain "stock-service/internal/domain/product_type"
)

type ListTypesInput struct {
	ProductID int32
}

type ListTypesOutput struct {
	Types []*producttypedomain.ProductType
}

type ListTypesUseCase struct {
	repo producttypedomain.Repository
}

func NewListTypesUseCase(repo producttypedomain.Repository) *ListTypesUseCase {
	return &ListTypesUseCase{repo: repo}
}

func (uc *ListTypesUseCase) Execute(input ListTypesInput) (*ListTypesOutput, error) {
	types, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	return &ListTypesOutput{Types: types}, nil
}
