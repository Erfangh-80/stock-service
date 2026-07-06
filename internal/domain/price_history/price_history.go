package pricehistory

import "time"

type PriceHistory struct {
	ID          int64
	ProductID   int32
	OldPrice    float64
	NewPrice    float64
	ChangedBy   string
	Description *string
	CreatedAt   time.Time
}

func NewPriceHistory(productID int32, oldPrice, newPrice float64, changedBy string, description *string) (*PriceHistory, error) {
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if oldPrice < 0 {
		return nil, ErrInvalidPrice
	}
	if newPrice < 0 {
		return nil, ErrInvalidPrice
	}
	if changedBy == "" {
		return nil, ErrChangedByRequired
	}

	return &PriceHistory{
		ProductID:   productID,
		OldPrice:    oldPrice,
		NewPrice:    newPrice,
		ChangedBy:   changedBy,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}
