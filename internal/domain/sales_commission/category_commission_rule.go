package salescommission

import "time"

type CategoryCommissionRule struct {
	ID          int64
	CategoryID  int32
	RatePercent float64
	MinPrice    float64
	MaxPrice    *float64
	IsActive    bool
	CreatedAt   time.Time
}

func NewCategoryCommissionRule(categoryID int32, ratePercent, minPrice float64) (*CategoryCommissionRule, error) {
	if err := ValidateRatePercent(ratePercent); err != nil {
		return nil, err
	}
	if err := ValidateMinPrice(minPrice); err != nil {
		return nil, err
	}
	return &CategoryCommissionRule{
		CategoryID:  categoryID,
		RatePercent: ratePercent,
		MinPrice:    minPrice,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}, nil
}

func (r *CategoryCommissionRule) UpdateMaxPrice(maxPrice float64) error {
	if err := ValidateMaxPrice(r.MinPrice, maxPrice); err != nil {
		return err
	}
	r.MaxPrice = &maxPrice
	return nil
}

func (r *CategoryCommissionRule) Activate() {
	r.IsActive = true
}

func (r *CategoryCommissionRule) Deactivate() {
	r.IsActive = false
}
