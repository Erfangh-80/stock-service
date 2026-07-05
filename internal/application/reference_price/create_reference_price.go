package referenceprice

import domainreferenceprice "stock-service/internal/domain/reference_price"

type CreateReferencePriceUseCase struct {
	repo domainreferenceprice.Repository
}

func NewCreateReferencePriceUseCase(repo domainreferenceprice.Repository) *CreateReferencePriceUseCase {
	return &CreateReferencePriceUseCase{repo: repo}
}

func (uc *CreateReferencePriceUseCase) Execute(productID int32, price float64, source string) (*domainreferenceprice.ReferencePrice, error) {
	rp, err := domainreferenceprice.NewReferencePrice(productID, price, source)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(rp); err != nil {
		return nil, err
	}
	return rp, nil
}
