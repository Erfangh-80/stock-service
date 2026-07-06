package promotion

import "time"

type DiscountType string

const (
	DiscountTypePercentage  DiscountType = "percentage"
	DiscountTypeFixedAmount DiscountType = "fixed_amount"
)

type PromotionStatus string

const (
	PromotionStatusInactive PromotionStatus = "inactive"
	PromotionStatusActive   PromotionStatus = "active"
)

type Promotion struct {
	ID                      int64
	Title                   string
	DiscountType            DiscountType
	DiscountValue           float64
	MinPurchase             *float64
	CouponCode              *string
	UsageLimit              *int
	UsedCount               int
	MaxDiscountAmount       *float64
	Budget                  *float64
	BudgetSpent             float64
	EligibleStoreIDs        []int64
	EligibleCategoryIDs     []int64
	EligibleProductIDs      []int32
	EligibleUserIDs         []int64
	RequiresApproval        bool
	StartAt                 *time.Time
	EndAt                   *time.Time
	IsCountdown             bool
	ExpireSaleWithPromotion bool
	Status                  PromotionStatus
	CreatedAt               time.Time
}

type CreatePromotionInput struct {
	Title                   string
	DiscountType            DiscountType
	DiscountValue           float64
	MinPurchase             *float64
	CouponCode              *string
	UsageLimit              *int
	MaxDiscountAmount       *float64
	Budget                  *float64
	EligibleStoreIDs        []int64
	EligibleCategoryIDs     []int64
	EligibleProductIDs      []int32
	EligibleUserIDs         []int64
	RequiresApproval        bool
	StartAt                 *time.Time
	EndAt                   *time.Time
	IsCountdown             bool
	ExpireSaleWithPromotion bool
}

type UpdatePromotionInput struct {
	Title                   *string
	DiscountType            *DiscountType
	DiscountValue           *float64
	MinPurchase             *float64
	CouponCode              *string
	UsageLimit              *int
	MaxDiscountAmount       *float64
	Budget                  *float64
	EligibleStoreIDs        []int64
	EligibleCategoryIDs     []int64
	EligibleProductIDs      []int32
	EligibleUserIDs         []int64
	RequiresApproval        *bool
	StartAt                 *time.Time
	EndAt                   *time.Time
	IsCountdown             *bool
	ExpireSaleWithPromotion *bool
}

func NewPromotion(input CreatePromotionInput) (*Promotion, error) {
	if err := ValidateTitle(input.Title); err != nil {
		return nil, err
	}
	if err := ValidateDiscountType(input.DiscountType); err != nil {
		return nil, err
	}
	if err := ValidateDiscountValue(input.DiscountType, input.DiscountValue); err != nil {
		return nil, err
	}
	if input.CouponCode != nil {
		if err := ValidateCouponCode(*input.CouponCode); err != nil {
			return nil, err
		}
	}
	if err := ValidatePromotionDates(input.StartAt, input.EndAt); err != nil {
		return nil, err
	}

	return &Promotion{
		Title:                   input.Title,
		DiscountType:            input.DiscountType,
		DiscountValue:           input.DiscountValue,
		MinPurchase:             input.MinPurchase,
		CouponCode:              input.CouponCode,
		UsageLimit:              input.UsageLimit,
		MaxDiscountAmount:       input.MaxDiscountAmount,
		Budget:                  input.Budget,
		EligibleStoreIDs:        input.EligibleStoreIDs,
		EligibleCategoryIDs:     input.EligibleCategoryIDs,
		EligibleProductIDs:      input.EligibleProductIDs,
		EligibleUserIDs:         input.EligibleUserIDs,
		RequiresApproval:        input.RequiresApproval,
		StartAt:                 input.StartAt,
		EndAt:                   input.EndAt,
		IsCountdown:             input.IsCountdown,
		ExpireSaleWithPromotion: input.ExpireSaleWithPromotion,
		Status:                  PromotionStatusInactive,
		CreatedAt:               time.Now(),
	}, nil
}

func (p *Promotion) Update(input UpdatePromotionInput) error {
	if input.Title != nil {
		if err := ValidateTitle(*input.Title); err != nil {
			return err
		}
		p.Title = *input.Title
	}
	if input.DiscountType != nil {
		if err := ValidateDiscountType(*input.DiscountType); err != nil {
			return err
		}
		p.DiscountType = *input.DiscountType
	}
	if input.DiscountValue != nil {
		if err := ValidateDiscountValue(p.DiscountType, *input.DiscountValue); err != nil {
			return err
		}
		p.DiscountValue = *input.DiscountValue
	}
	if input.MinPurchase != nil {
		p.MinPurchase = input.MinPurchase
	}
	if input.CouponCode != nil {
		if err := ValidateCouponCode(*input.CouponCode); err != nil {
			return err
		}
		p.CouponCode = input.CouponCode
	}
	if input.UsageLimit != nil {
		p.UsageLimit = input.UsageLimit
	}
	if input.MaxDiscountAmount != nil {
		p.MaxDiscountAmount = input.MaxDiscountAmount
	}
	if input.Budget != nil {
		p.Budget = input.Budget
	}
	if input.EligibleStoreIDs != nil {
		p.EligibleStoreIDs = input.EligibleStoreIDs
	}
	if input.EligibleCategoryIDs != nil {
		p.EligibleCategoryIDs = input.EligibleCategoryIDs
	}
	if input.EligibleProductIDs != nil {
		p.EligibleProductIDs = input.EligibleProductIDs
	}
	if input.EligibleUserIDs != nil {
		p.EligibleUserIDs = input.EligibleUserIDs
	}
	if input.RequiresApproval != nil {
		p.RequiresApproval = *input.RequiresApproval
	}

	startAt := p.StartAt
	endAt := p.EndAt
	if input.StartAt != nil {
		startAt = input.StartAt
	}
	if input.EndAt != nil {
		endAt = input.EndAt
	}
	if input.StartAt != nil || input.EndAt != nil {
		if err := ValidatePromotionDates(startAt, endAt); err != nil {
			return err
		}
		p.StartAt = startAt
		p.EndAt = endAt
	}

	if input.IsCountdown != nil {
		p.IsCountdown = *input.IsCountdown
	}
	if input.ExpireSaleWithPromotion != nil {
		p.ExpireSaleWithPromotion = *input.ExpireSaleWithPromotion
	}

	return nil
}

func (p *Promotion) Activate() {
	p.Status = PromotionStatusActive
}

func (p *Promotion) Deactivate() {
	p.Status = PromotionStatusInactive
}

func (p *Promotion) IsActive() bool {
	return p.Status == PromotionStatusActive
}

func (p *Promotion) IsExpired() bool {
	if p.EndAt == nil {
		return false
	}
	return time.Now().After(*p.EndAt)
}

func (p *Promotion) IsScheduled() bool {
	if p.StartAt == nil {
		return false
	}
	return time.Now().Before(*p.StartAt)
}

func (p *Promotion) CanApply() error {
	if !p.IsActive() {
		return ErrPromotionNotActive
	}
	if p.IsExpired() {
		return ErrPromotionExpired
	}
	if p.UsageLimit != nil && p.UsedCount >= *p.UsageLimit {
		return ErrPromotionUsageLimitExceeded
	}
	if p.Budget != nil && p.BudgetSpent >= *p.Budget {
		return ErrPromotionBudgetExceeded
	}
	return nil
}

func (p *Promotion) RecordUsage() {
	p.UsedCount++
}

func (p *Promotion) SpendBudget(amount float64) {
	p.BudgetSpent += amount
}

func (p *Promotion) CalculateDiscountPrice(basePrice float64) float64 {
	switch p.DiscountType {
	case DiscountTypePercentage:
		discount := basePrice * p.DiscountValue / 100
		if p.MaxDiscountAmount != nil && discount > *p.MaxDiscountAmount {
			discount = *p.MaxDiscountAmount
		}
		if p.MinPurchase != nil && basePrice < *p.MinPurchase {
			return basePrice
		}
		result := basePrice - discount
		if result < 0 {
			return 0
		}
		return result
	case DiscountTypeFixedAmount:
		if p.MinPurchase != nil && basePrice < *p.MinPurchase {
			return basePrice
		}
		result := basePrice - p.DiscountValue
		if result < 0 {
			return 0
		}
		return result
	default:
		return basePrice
	}
}

func (p *Promotion) IsEligibleForStore(storeID int64) bool {
	if len(p.EligibleStoreIDs) == 0 {
		return true
	}
	for _, id := range p.EligibleStoreIDs {
		if id == storeID {
			return true
		}
	}
	return false
}

func (p *Promotion) IsEligibleForCategory(categoryID int64) bool {
	if len(p.EligibleCategoryIDs) == 0 {
		return true
	}
	for _, id := range p.EligibleCategoryIDs {
		if id == categoryID {
			return true
		}
	}
	return false
}

func (p *Promotion) IsEligibleForProduct(productID int32) bool {
	if len(p.EligibleProductIDs) == 0 {
		return true
	}
	for _, id := range p.EligibleProductIDs {
		if id == productID {
			return true
		}
	}
	return false
}

func (p *Promotion) IsEligibleForUser(userID int64) bool {
	if len(p.EligibleUserIDs) == 0 {
		return true
	}
	for _, id := range p.EligibleUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}
