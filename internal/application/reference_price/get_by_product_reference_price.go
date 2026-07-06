package referenceprice

import domainreferenceprice "stock-service/internal/domain/reference_price"

type GetByProductReferencePriceInput struct {
	ProductID int32
}

type GetByProductReferencePriceUseCase struct {
	repo domainreferenceprice.Repository
}

func NewGetByProductReferencePriceUseCase(repo domainreferenceprice.Repository) *GetByProductReferencePriceUseCase {
	return &GetByProductReferencePriceUseCase{repo: repo}
}

func (uc *GetByProductReferencePriceUseCase) Execute(input GetByProductReferencePriceInput) (*domainreferenceprice.ReferencePrice, error) {
	rp, err := uc.repo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	if rp == nil {
		return nil, domainreferenceprice.ErrReferencePriceNotFound
	}
	return rp, nil
}
