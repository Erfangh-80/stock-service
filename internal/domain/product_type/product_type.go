package producttype

import "time"

type ProductType struct {
	ID        int64
	ProductID int32
	Name      string
	Value     string
	SortOrder int
	CreatedAt time.Time
}

func NewProductType(productID int32, name, value string, sortOrder int) (*ProductType, error) {
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if name == "" {
		return nil, ErrNameRequired
	}
	if value == "" {
		return nil, ErrValueRequired
	}

	return &ProductType{
		ProductID: productID,
		Name:      name,
		Value:     value,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}, nil
}

func (pt *ProductType) UpdateValue(value string) error {
	if value == "" {
		return ErrValueRequired
	}
	pt.Value = value
	return nil
}

func (pt *ProductType) UpdateSortOrder(order int) {
	pt.SortOrder = order
}
