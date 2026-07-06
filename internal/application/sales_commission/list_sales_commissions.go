package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type ListSalesCommissionsInput struct {
	InventoryID *int64
	SaleModel   *string
	Page        int
	Limit       int
}

type ListSalesCommissionsOutput struct {
	Commissions []*domainsalescommission.SalesCommission
	Total       int
	Page        int
	Limit       int
}

type ListSalesCommissionsUseCase struct {
	repo domainsalescommission.Repository
}

func NewListSalesCommissionsUseCase(repo domainsalescommission.Repository) *ListSalesCommissionsUseCase {
	return &ListSalesCommissionsUseCase{repo: repo}
}

func (uc *ListSalesCommissionsUseCase) Execute(input ListSalesCommissionsInput) (*ListSalesCommissionsOutput, error) {
	filter := domainsalescommission.SalesCommissionFilter{
		InventoryID: input.InventoryID,
		SaleModel:   input.SaleModel,
		Page:        input.Page,
		Limit:       input.Limit,
	}
	items, total, err := uc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}
	return &ListSalesCommissionsOutput{
		Commissions: items,
		Total:       total,
		Page:        input.Page,
		Limit:       input.Limit,
	}, nil
}
