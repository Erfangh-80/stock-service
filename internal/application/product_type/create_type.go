package producttype

import (
	producttypedomain "stock-service/internal/domain/product_type"
)

type CreateTypeInput struct {
	ProductID int32
	Name      string
	Value     string
	SortOrder int
}

type CreateTypeUseCase struct {
	repo producttypedomain.Repository
}

func NewCreateTypeUseCase(repo producttypedomain.Repository) *CreateTypeUseCase {
	return &CreateTypeUseCase{repo: repo}
}

func (uc *CreateTypeUseCase) Execute(input CreateTypeInput) (*producttypedomain.ProductType, error) {
	pt, err := producttypedomain.NewProductType(input.ProductID, input.Name, input.Value, input.SortOrder)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(pt); err != nil {
		return nil, err
	}
	return pt, nil
}
