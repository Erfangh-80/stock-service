package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type GetSalesCommissionInput struct {
	ID int64
}

type GetSalesCommissionUseCase struct {
	repo domainsalescommission.Repository
}

func NewGetSalesCommissionUseCase(repo domainsalescommission.Repository) *GetSalesCommissionUseCase {
	return &GetSalesCommissionUseCase{repo: repo}
}

func (uc *GetSalesCommissionUseCase) Execute(input GetSalesCommissionInput) (*domainsalescommission.SalesCommission, error) {
	sc, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, domainsalescommission.ErrCommissionNotFound
	}
	return sc, nil
}
