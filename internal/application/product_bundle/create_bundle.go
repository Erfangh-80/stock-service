package productbundle

import (
	productbundledomain "stock-service/internal/domain/product_bundle"
)

type CreateBundleInput struct {
	ProductID        int32
	RelatedProductID int32
	Type             string
	SortOrder        int
}

type CreateBundleUseCase struct {
	repo productbundledomain.Repository
}

func NewCreateBundleUseCase(repo productbundledomain.Repository) *CreateBundleUseCase {
	return &CreateBundleUseCase{repo: repo}
}

func (uc *CreateBundleUseCase) Execute(input CreateBundleInput) (*productbundledomain.ProductBundle, error) {
	pb, err := productbundledomain.NewProductBundle(
		input.ProductID, input.RelatedProductID,
		productbundledomain.BundleType(input.Type), input.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(pb); err != nil {
		return nil, err
	}
	return pb, nil
}
