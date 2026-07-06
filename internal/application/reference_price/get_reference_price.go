package referenceprice

import domainreferenceprice "stock-service/internal/domain/reference_price"

type GetReferencePriceInput struct {
	ID int64
}

type GetReferencePriceUseCase struct {
	repo domainreferenceprice.Repository
}

func NewGetReferencePriceUseCase(repo domainreferenceprice.Repository) *GetReferencePriceUseCase {
	return &GetReferencePriceUseCase{repo: repo}
}

func (uc *GetReferencePriceUseCase) Execute(input GetReferencePriceInput) (*domainreferenceprice.ReferencePrice, error) {
	rp, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if rp == nil {
		return nil, domainreferenceprice.ErrReferencePriceNotFound
	}
	return rp, nil
}
