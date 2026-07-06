package promotioninterface

import (
	"time"

	apppromotion "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
	iface "stock-service/internal/interface"
)

type PromotionResponse struct {
	ID                      int64               `json:"id"`
	Title                   string              `json:"title"`
	DiscountType            string              `json:"discount_type"`
	DiscountValue           float64             `json:"discount_value"`
	MinPurchase             *float64            `json:"min_purchase,omitempty"`
	CouponCode              *string             `json:"coupon_code,omitempty"`
	UsageLimit              *int                `json:"usage_limit,omitempty"`
	UsedCount               int                 `json:"used_count"`
	MaxDiscountAmount       *float64            `json:"max_discount_amount,omitempty"`
	Budget                  *float64            `json:"budget,omitempty"`
	BudgetSpent             float64             `json:"budget_spent"`
	EligibleStoreIDs        []int64             `json:"eligible_store_ids,omitempty"`
	EligibleCategoryIDs     []int64             `json:"eligible_category_ids,omitempty"`
	EligibleProductIDs      []int32             `json:"eligible_product_ids,omitempty"`
	EligibleUserIDs         []int64             `json:"eligible_user_ids,omitempty"`
	RequiresApproval        bool                `json:"requires_approval"`
	StartAt                 *string             `json:"start_at,omitempty"`
	EndAt                   *string             `json:"end_at,omitempty"`
	IsCountdown             bool                `json:"is_countdown"`
	ExpireSaleWithPromotion bool                `json:"expire_sale_with_promotion"`
	Status                  string              `json:"status"`
	CreatedAt               string              `json:"created_at"`
}

type PromotionListItem struct {
	ID            int64   `json:"id"`
	Title         string  `json:"title"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	Status        string  `json:"status"`
	StartAt       *string `json:"start_at,omitempty"`
	EndAt         *string `json:"end_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type PromotionListResponse struct {
	Promotions []PromotionListItem `json:"promotions"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}

type CreatePromotionParams struct {
	Title                   string   `json:"title"`
	DiscountType            string   `json:"discount_type"`
	DiscountValue           float64  `json:"discount_value"`
	MinPurchase             *float64 `json:"min_purchase,omitempty"`
	CouponCode              *string  `json:"coupon_code,omitempty"`
	UsageLimit              *int     `json:"usage_limit,omitempty"`
	MaxDiscountAmount       *float64 `json:"max_discount_amount,omitempty"`
	Budget                  *float64 `json:"budget,omitempty"`
	EligibleStoreIDs        []int64  `json:"eligible_store_ids,omitempty"`
	EligibleCategoryIDs     []int64  `json:"eligible_category_ids,omitempty"`
	EligibleProductIDs      []int32  `json:"eligible_product_ids,omitempty"`
	EligibleUserIDs         []int64  `json:"eligible_user_ids,omitempty"`
	RequiresApproval        bool     `json:"requires_approval"`
	StartAt                 *string  `json:"start_at,omitempty"`
	EndAt                   *string  `json:"end_at,omitempty"`
	IsCountdown             bool     `json:"is_countdown"`
	ExpireSaleWithPromotion bool     `json:"expire_sale_with_promotion"`
}

type UpdatePromotionParams struct {
	ID                      int64
	Title                   *string  `json:"title,omitempty"`
	DiscountType            *string  `json:"discount_type,omitempty"`
	DiscountValue           *float64 `json:"discount_value,omitempty"`
	MinPurchase             *float64 `json:"min_purchase,omitempty"`
	CouponCode              *string  `json:"coupon_code,omitempty"`
	UsageLimit              *int     `json:"usage_limit,omitempty"`
	MaxDiscountAmount       *float64 `json:"max_discount_amount,omitempty"`
	Budget                  *float64 `json:"budget,omitempty"`
	EligibleStoreIDs        []int64  `json:"eligible_store_ids,omitempty"`
	EligibleCategoryIDs     []int64  `json:"eligible_category_ids,omitempty"`
	EligibleProductIDs      []int32  `json:"eligible_product_ids,omitempty"`
	EligibleUserIDs         []int64  `json:"eligible_user_ids,omitempty"`
	RequiresApproval        *bool    `json:"requires_approval,omitempty"`
	StartAt                 *string  `json:"start_at,omitempty"`
	EndAt                   *string  `json:"end_at,omitempty"`
	IsCountdown             *bool    `json:"is_countdown,omitempty"`
	ExpireSaleWithPromotion *bool    `json:"expire_sale_with_promotion,omitempty"`
}

type ListPromotionsParams struct {
	Status       *string
	DiscountType *string
	Search       *string
	Page         int
	Limit        int
}

type createPromotionUseCase interface {
	Execute(input apppromotion.CreatePromotionInput) (*promotion.Promotion, error)
}

type getPromotionUseCase interface {
	Execute(input apppromotion.GetPromotionInput) (*promotion.Promotion, error)
}

type activatePromotionUseCase interface {
	Execute(id int64) error
}

type deactivatePromotionUseCase interface {
	Execute(id int64) error
}

type updatePromotionUseCase interface {
	Execute(input apppromotion.UpdatePromotionInput) (*promotion.Promotion, error)
}

type deletePromotionUseCase interface {
	Execute(input apppromotion.DeletePromotionInput) error
}

type listPromotionsUseCase interface {
	Execute(input apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error)
}

type Adapter struct {
	create     createPromotionUseCase
	get        getPromotionUseCase
	activate   activatePromotionUseCase
	deactivate deactivatePromotionUseCase
	update     updatePromotionUseCase
	deleteUC   deletePromotionUseCase
	list       listPromotionsUseCase
}

func NewAdapter(
	create createPromotionUseCase,
	get getPromotionUseCase,
	activate activatePromotionUseCase,
	deactivate deactivatePromotionUseCase,
	update updatePromotionUseCase,
	deleteUC deletePromotionUseCase,
	list listPromotionsUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, activate: activate, deactivate: deactivate,
		update: update, deleteUC: deleteUC, list: list,
	}
}

func (a *Adapter) Create(params CreatePromotionParams) (*PromotionResponse, error) {
	dt := promotion.DiscountType(params.DiscountType)
	input := apppromotion.CreatePromotionInput{
		Title:                   params.Title,
		DiscountType:            dt,
		DiscountValue:           params.DiscountValue,
		MinPurchase:             params.MinPurchase,
		CouponCode:              params.CouponCode,
		UsageLimit:              params.UsageLimit,
		MaxDiscountAmount:       params.MaxDiscountAmount,
		Budget:                  params.Budget,
		EligibleStoreIDs:        params.EligibleStoreIDs,
		EligibleCategoryIDs:     params.EligibleCategoryIDs,
		EligibleProductIDs:      params.EligibleProductIDs,
		EligibleUserIDs:         params.EligibleUserIDs,
		RequiresApproval:        params.RequiresApproval,
		StartAt:                 params.StartAt,
		EndAt:                   params.EndAt,
		IsCountdown:             params.IsCountdown,
		ExpireSaleWithPromotion: params.ExpireSaleWithPromotion,
	}
	result, err := a.create.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*PromotionResponse, error) {
	result, err := a.get.Execute(apppromotion.GetPromotionInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Update(id int64, params UpdatePromotionParams) (*PromotionResponse, error) {
	input := apppromotion.UpdatePromotionInput{
		ID:                      id,
		Title:                   params.Title,
		MinPurchase:             params.MinPurchase,
		CouponCode:              params.CouponCode,
		UsageLimit:              params.UsageLimit,
		MaxDiscountAmount:       params.MaxDiscountAmount,
		Budget:                  params.Budget,
		EligibleStoreIDs:        params.EligibleStoreIDs,
		EligibleCategoryIDs:     params.EligibleCategoryIDs,
		EligibleProductIDs:      params.EligibleProductIDs,
		EligibleUserIDs:         params.EligibleUserIDs,
		RequiresApproval:        params.RequiresApproval,
		StartAt:                 params.StartAt,
		EndAt:                   params.EndAt,
		IsCountdown:             params.IsCountdown,
		ExpireSaleWithPromotion: params.ExpireSaleWithPromotion,
	}
	if params.DiscountType != nil {
		dt := promotion.DiscountType(*params.DiscountType)
		input.DiscountType = &dt
	}
	if params.DiscountValue != nil {
		input.DiscountValue = params.DiscountValue
	}
	result, err := a.update.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Activate(id int64) error {
	err := a.activate.Execute(id)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Deactivate(id int64) error {
	err := a.deactivate.Execute(id)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.deleteUC.Execute(apppromotion.DeletePromotionInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (a *Adapter) List(params ListPromotionsParams) (*PromotionListResponse, error) {
	input := apppromotion.ListPromotionsInput{
		Page:  params.Page,
		Limit: params.Limit,
		Search: params.Search,
	}
	if params.Status != nil {
		s := promotion.PromotionStatus(*params.Status)
		input.Status = &s
	}
	if params.DiscountType != nil {
		dt := promotion.DiscountType(*params.DiscountType)
		input.DiscountType = &dt
	}
	result, err := a.list.Execute(input)
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]PromotionListItem, 0, len(result.Promotions))
	for _, p := range result.Promotions {
		items = append(items, PromotionListItem{
			ID:            p.ID,
			Title:         p.Title,
			DiscountType:  string(p.DiscountType),
			DiscountValue: p.DiscountValue,
			Status:        string(p.Status),
			StartAt:       formatTime(p.StartAt),
			EndAt:         formatTime(p.EndAt),
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return &PromotionListResponse{
		Promotions: items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
	}, nil
}

func mapError(err error) error {
	switch err {
	case promotion.ErrPromotionNotFound:
		return iface.ErrNotFound
	case promotion.ErrTitleRequired, promotion.ErrTitleTooLong,
		promotion.ErrDiscountTypeRequired, promotion.ErrInvalidDiscountType,
		promotion.ErrDiscountValueRequired, promotion.ErrInvalidDiscountValue,
		promotion.ErrDiscountValueTooHigh,
		promotion.ErrInvalidPromotionDates, promotion.ErrInvalidCouponCode,
		promotion.ErrCouponCodeAlreadyExists, promotion.ErrInvalidDateFormat:
		return iface.ErrInvalidInput
	case promotion.ErrPromotionNotActive, promotion.ErrPromotionExpired,
		promotion.ErrPromotionNotStarted, promotion.ErrPromotionUsageLimitExceeded,
		promotion.ErrPromotionBudgetExceeded,
		promotion.ErrIneligibleStore, promotion.ErrIneligibleCategory,
		promotion.ErrIneligibleProduct, promotion.ErrIneligibleUser:
		return iface.ErrConflict
	default:
		return iface.ErrInternal
	}
}

func toResponse(p *promotion.Promotion) *PromotionResponse {
	return &PromotionResponse{
		ID:                      p.ID,
		Title:                   p.Title,
		DiscountType:            string(p.DiscountType),
		DiscountValue:           p.DiscountValue,
		MinPurchase:             p.MinPurchase,
		CouponCode:              p.CouponCode,
		UsageLimit:              p.UsageLimit,
		UsedCount:               p.UsedCount,
		MaxDiscountAmount:       p.MaxDiscountAmount,
		Budget:                  p.Budget,
		BudgetSpent:             p.BudgetSpent,
		EligibleStoreIDs:        p.EligibleStoreIDs,
		EligibleCategoryIDs:     p.EligibleCategoryIDs,
		EligibleProductIDs:      p.EligibleProductIDs,
		EligibleUserIDs:         p.EligibleUserIDs,
		RequiresApproval:        p.RequiresApproval,
		StartAt:                 formatTime(p.StartAt),
		EndAt:                   formatTime(p.EndAt),
		IsCountdown:             p.IsCountdown,
		ExpireSaleWithPromotion: p.ExpireSaleWithPromotion,
		Status:                  string(p.Status),
		CreatedAt:               p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02T15:04:05Z")
	return &s
}
