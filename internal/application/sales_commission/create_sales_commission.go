package salescommission

import (
	"stock-service/internal/domain/sales_commission"
)

type CreateSalesCommissionUseCase struct {
	repo salescommission.Repository
}

func NewCreateSalesCommissionUseCase(repo salescommission.Repository) *CreateSalesCommissionUseCase {
	return &CreateSalesCommissionUseCase{repo: repo}
}

func (uc *CreateSalesCommissionUseCase) Execute(inventoryID, categoryCommissionRuleID int64, saleModel salescommission.SaleModel, ratePercent, minPrice float64) (*salescommission.SalesCommission, error) {
	sc, err := salescommission.NewSalesCommission(inventoryID, categoryCommissionRuleID, saleModel, ratePercent, minPrice)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(sc); err != nil {
		return nil, err
	}
	return sc, nil
}
