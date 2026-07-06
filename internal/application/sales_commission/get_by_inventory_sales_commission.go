package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
)

type GetByInventorySalesCommissionInput struct {
	InventoryID int64
}

type GetByInventorySalesCommissionUseCase struct {
	repo domainsalescommission.Repository
}

func NewGetByInventorySalesCommissionUseCase(repo domainsalescommission.Repository) *GetByInventorySalesCommissionUseCase {
	return &GetByInventorySalesCommissionUseCase{repo: repo}
}

func (uc *GetByInventorySalesCommissionUseCase) Execute(input GetByInventorySalesCommissionInput) (*domainsalescommission.SalesCommission, error) {
	sc, err := uc.repo.FindByInventoryID(input.InventoryID)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, domainsalescommission.ErrCommissionNotFound
	}
	return sc, nil
}
