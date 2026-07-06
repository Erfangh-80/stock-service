package promotion

import (
	"time"

	domainpromotion "stock-service/internal/domain/promotion"
)

type UpdatePromotionInput struct {
	ID                      int64
	Title                   *string
	DiscountType            *domainpromotion.DiscountType
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
	StartAt                 *string
	EndAt                   *string
	IsCountdown             *bool
	ExpireSaleWithPromotion *bool
}

type UpdatePromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewUpdatePromotionUseCase(repo domainpromotion.Repository) *UpdatePromotionUseCase {
	return &UpdatePromotionUseCase{repo: repo}
}

func (uc *UpdatePromotionUseCase) Execute(input UpdatePromotionInput) (*domainpromotion.Promotion, error) {
	p, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, domainpromotion.ErrPromotionNotFound
	}

	if input.CouponCode != nil {
		existing, _ := uc.repo.FindByCouponCode(*input.CouponCode)
		if existing != nil && existing.ID != input.ID {
			return nil, domainpromotion.ErrCouponCodeAlreadyExists
		}
	}

	var startAt, endAt *time.Time
	if input.StartAt != nil {
		t, err := time.Parse(time.RFC3339, *input.StartAt)
		if err != nil {
			return nil, domainpromotion.ErrInvalidDateFormat
		}
		startAt = &t
	}
	if input.EndAt != nil {
		t, err := time.Parse(time.RFC3339, *input.EndAt)
		if err != nil {
			return nil, domainpromotion.ErrInvalidDateFormat
		}
		endAt = &t
	}

	domainInput := domainpromotion.UpdatePromotionInput{
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
		StartAt:                 startAt,
		EndAt:                   endAt,
		IsCountdown:             input.IsCountdown,
		ExpireSaleWithPromotion: input.ExpireSaleWithPromotion,
	}

	if err := p.Update(domainInput); err != nil {
		return nil, err
	}
	if err := uc.repo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}
