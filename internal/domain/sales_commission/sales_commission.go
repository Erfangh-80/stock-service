package salescommission

import "time"

// TODO: SaleModel is defined here and also in inventory. Consider a shared kernel
// if these are the same domain concept. Keeping separate for now to avoid
// cross-package domain dependency.
type SaleModel string

const SaleModelRetail SaleModel = "retail"

type SalesCommission struct {
	ID                       int64
	InventoryID              int64
	CategoryCommissionRuleID int64
	SaleModel                SaleModel
	RatePercent              float64
	MinQty                   *int
	MinPrice                 float64
	MaxPrice                 *float64
	CreatedAt                time.Time
}

func NewSalesCommission(inventoryID, categoryCommissionRuleID int64, saleModel SaleModel, ratePercent, minPrice float64) (*SalesCommission, error) {
	if err := ValidateRatePercent(ratePercent); err != nil {
		return nil, err
	}
	if err := ValidateMinPrice(minPrice); err != nil {
		return nil, err
	}

	return &SalesCommission{
		InventoryID:              inventoryID,
		CategoryCommissionRuleID: categoryCommissionRuleID,
		SaleModel:                saleModel,
		RatePercent:              ratePercent,
		MinPrice:                 minPrice,
		CreatedAt:                time.Now(),
	}, nil
}

func (sc *SalesCommission) UpdateMaxPrice(maxPrice float64) error {
	if err := ValidateMaxPrice(sc.MinPrice, maxPrice); err != nil {
		return err
	}
	sc.MaxPrice = &maxPrice
	return nil
}

func (sc *SalesCommission) UpdateMinQty(qty int) error {
	if qty < 0 {
		return ErrInvalidMinQty
	}
	sc.MinQty = &qty
	return nil
}
