package dto

type CreateReferencePriceRequest struct {
	ProductID int32   `json:"product_id"`
	Price     float64 `json:"price"`
	Source    string  `json:"source"`
}
