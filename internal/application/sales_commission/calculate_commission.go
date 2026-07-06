package salescommission

import (
	domainsalescommission "stock-service/internal/domain/sales_commission"
	inventorydomain "stock-service/internal/domain/inventory"
)

type CalculateCommissionInput struct {
	InventoryID int64
	Quantity    int
}

type CommissionCalculation struct {
	CommissionID   int64
	InventoryID    int64
	BasePriceUsed  float64
	Quantity       int
	RatePercent    float64
	CommissionAmt  float64
	PriceSource    string
}

type CalculateCommissionUseCase struct {
	commRepo  domainsalescommission.Repository
	invRepo   inventorydomain.Repository
}

func NewCalculateCommissionUseCase(commRepo domainsalescommission.Repository, invRepo inventorydomain.Repository) *CalculateCommissionUseCase {
	return &CalculateCommissionUseCase{commRepo: commRepo, invRepo: invRepo}
}

func (uc *CalculateCommissionUseCase) Execute(input CalculateCommissionInput) (*CommissionCalculation, error) {
	sc, err := uc.commRepo.FindByInventoryID(input.InventoryID)
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, domainsalescommission.ErrCommissionNotFound
	}

	inv, err := uc.invRepo.FindByID(input.InventoryID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, inventorydomain.ErrInventoryNotFound
	}

	priceSource := "base_price"
	price := inv.BasePrice
	if sc.MinPrice > 0 && price < sc.MinPrice {
		price = sc.MinPrice
	}
	if sc.MaxPrice != nil && price > *sc.MaxPrice {
		price = *sc.MaxPrice
	}

	if inv.FinalPrice != nil && inv.PromotionID != nil {
		priceSource = "final_price"
		price = *inv.FinalPrice
		if sc.MinPrice > 0 && price < sc.MinPrice {
			price = sc.MinPrice
		}
		if sc.MaxPrice != nil && price > *sc.MaxPrice {
			price = *sc.MaxPrice
		}
	}

	qty := input.Quantity
	if sc.MinQty != nil && qty < *sc.MinQty {
		qty = *sc.MinQty
	}

	commissionAmt := price * sc.RatePercent / 100 * float64(qty)

	return &CommissionCalculation{
		CommissionID:  sc.ID,
		InventoryID:   input.InventoryID,
		BasePriceUsed: price,
		Quantity:      qty,
		RatePercent:   sc.RatePercent,
		CommissionAmt: commissionAmt,
		PriceSource:   priceSource,
	}, nil
}
