package dto

import "time"

type CreateInventoryRequest struct {
	StoreID     int64   `json:"store_id"`
	WarehouseID int64   `json:"warehouse_id"`
	ProductID   int32   `json:"product_id"`
	BasePrice   float64 `json:"base_price"`
}

type ApplyPromotionRequest struct {
	PromotionID int64     `json:"promotion_id"`
	FinalPrice  float64   `json:"final_price"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

type UpdateInventoryRequest struct {
	InstantQty   int            `json:"instant_qty"`
	ScheduledQty map[string]int `json:"scheduled_qty"`
	MinOrderQty  int            `json:"min_order_qty"`
	MaxOrderQty  *int           `json:"max_order_qty"`
}

type ListInventoryQuery struct {
	StoreID          *int64  `json:"store_id"`
	ProductID        *int32  `json:"product_id"`
	VendorSaleStatus *string `json:"vendor_sale_status"`
	SystemSaleStatus *string `json:"system_sale_status"`
	Page             int     `json:"page"`
	Limit            int     `json:"limit"`
}

type SearchInventoryQuery struct {
	Query string `json:"query"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type ReserveQuantityRequest struct {
	Quantity int `json:"quantity"`
}

type ReleaseQuantityRequest struct {
	Quantity int `json:"quantity"`
}

type CheckLowStockRequest struct {
	Threshold int `json:"threshold"`
}
