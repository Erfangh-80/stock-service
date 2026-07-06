package productattribute

import "time"

type ProductAttribute struct {
	ID        int64
	ProductID int32
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProductAttribute(productID int32, key, value string) (*ProductAttribute, error) {
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if key == "" {
		return nil, ErrKeyRequired
	}

	now := time.Now()
	return &ProductAttribute{
		ProductID: productID,
		Key:       key,
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (pa *ProductAttribute) UpdateValue(value string) {
	pa.Value = value
	pa.UpdatedAt = time.Now()
}
