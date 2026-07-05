package referenceprice

import "time"

type ReferencePrice struct {
	ID        int64
	ProductID int32
	Price     float64
	Source    string
	CreatedAt time.Time
}

func NewReferencePrice(productID int32, price float64, source string) (*ReferencePrice, error) {
	if err := ValidatePrice(price); err != nil {
		return nil, err
	}

	return &ReferencePrice{
		ProductID: productID,
		Price:     price,
		Source:    source,
		CreatedAt: time.Now(),
	}, nil
}
