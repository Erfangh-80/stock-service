package inventory

import (
	"time"

	"stock-service/internal/domain/inventory"
	promotiondomain "stock-service/internal/domain/promotion"
)

type ApplyPromotionInput struct {
	SaleID      int64
	PromotionID int64
	FinalPrice  float64
	StartAt     time.Time
	EndAt       time.Time
}

type ApplyPromotionUseCase interface {
	Execute(input ApplyPromotionInput) (*inventory.Inventory, error)
}

type applyPromotionInteractor struct {
	repo        inventory.Repository
	promotionRepo promotiondomain.Repository
}

func NewApplyPromotionUseCase(repo inventory.Repository, promotionRepo promotiondomain.Repository) ApplyPromotionUseCase {
	return &applyPromotionInteractor{repo: repo, promotionRepo: promotionRepo}
}

func (uc *applyPromotionInteractor) Execute(input ApplyPromotionInput) (*inventory.Inventory, error) {
	inv, err := uc.repo.FindByID(input.SaleID)
	if err != nil {
		return nil, err
	}

	promo, err := uc.promotionRepo.FindByID(input.PromotionID)
	if err != nil {
		return nil, err
	}
	if promo == nil {
		return nil, promotiondomain.ErrPromotionNotFound
	}

	if err := promo.CanApply(); err != nil {
		return nil, err
	}

	if !promo.IsEligibleForStore(inv.StoreID) {
		return nil, promotiondomain.ErrIneligibleStore
	}

	if !promo.IsEligibleForProduct(inv.ProductID) {
		return nil, promotiondomain.ErrIneligibleProduct
	}

	finalPrice := input.FinalPrice
	if finalPrice == 0 {
		finalPrice = promo.CalculateDiscountPrice(inv.BasePrice)
	}

	startAt := input.StartAt
	endAt := input.EndAt
	if startAt.IsZero() && promo.StartAt != nil {
		startAt = *promo.StartAt
	}
	if endAt.IsZero() && promo.EndAt != nil {
		endAt = *promo.EndAt
	}

	if err := inv.ApplyPromotion(input.PromotionID, finalPrice, startAt, endAt); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(inv); err != nil {
		return nil, err
	}

	promo.RecordUsage()
	discount := inv.BasePrice - finalPrice
	if discount > 0 {
		promo.SpendBudget(discount)
	}
	if err := uc.promotionRepo.Save(promo); err != nil {
		return nil, err
	}

	return inv, nil
}
