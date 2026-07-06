package referenceprice

import domainreferenceprice "stock-service/internal/domain/reference_price"

type ListReferencePricesInput struct {
	ProductID *int32
	Source    *string
	Page      int
	Limit     int
}

type ListReferencePricesOutput struct {
	ReferencePrices []*domainreferenceprice.ReferencePrice
	Total           int
	Page            int
	Limit           int
}

type ListReferencePricesUseCase struct {
	repo domainreferenceprice.Repository
}

func NewListReferencePricesUseCase(repo domainreferenceprice.Repository) *ListReferencePricesUseCase {
	return &ListReferencePricesUseCase{repo: repo}
}

func (uc *ListReferencePricesUseCase) Execute(input ListReferencePricesInput) (*ListReferencePricesOutput, error) {
	filter := domainreferenceprice.ReferencePriceFilter{
		ProductID: input.ProductID,
		Source:    input.Source,
		Page:      input.Page,
		Limit:     input.Limit,
	}
	items, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListReferencePricesOutput{
		ReferencePrices: items,
		Total:           total,
		Page:            input.Page,
		Limit:           input.Limit,
	}, nil
}
