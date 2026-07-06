package productbundle

import (
	productbundledomain "stock-service/internal/domain/product_bundle"
)

type ListBundlesInput struct {
	ProductID int32
}

type ListBundlesOutput struct {
	Bundles []*productbundledomain.ProductBundle
}

type ListBundlesUseCase struct {
	repo productbundledomain.Repository
}

func NewListBundlesUseCase(repo productbundledomain.Repository) *ListBundlesUseCase {
	return &ListBundlesUseCase{repo: repo}
}

func (uc *ListBundlesUseCase) Execute(input ListBundlesInput) (*ListBundlesOutput, error) {
	bundles, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	return &ListBundlesOutput{Bundles: bundles}, nil
}
