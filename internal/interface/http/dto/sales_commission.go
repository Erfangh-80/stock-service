package dto

type CreateSalesCommissionRequest struct {
	InventoryID              int64   `json:"inventory_id"`
	CategoryCommissionRuleID int64   `json:"category_commission_rule_id"`
	SaleModel                string  `json:"sale_model"`
	RatePercent              float64 `json:"rate_percent"`
	MinPrice                 float64 `json:"min_price"`
}

type UpdateMaxPriceRequest struct {
	MaxPrice float64 `json:"max_price"`
}

type UpdateMinQtyRequest struct {
	MinQty int `json:"min_qty"`
}

type CalculateCommissionRequest struct {
	InventoryID int64 `json:"inventory_id"`
	Quantity    int   `json:"quantity"`
}

type CreateCategoryCommissionRuleRequest struct {
	CategoryID  int32   `json:"category_id"`
	RatePercent float64 `json:"rate_percent"`
	MinPrice    float64 `json:"min_price"`
}

type UpdateCategoryCommissionRuleRequest struct {
	RatePercent *float64 `json:"rate_percent,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	Activate    *bool    `json:"activate,omitempty"`
}
