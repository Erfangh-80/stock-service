package inventory

import (
	"testing"
	"time"

	"stock-service/internal/domain/inventory"
)

func TestNewInventory_Success(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if inv.StoreID != 1 {
		t.Errorf("expected StoreID 1, got %d", inv.StoreID)
	}
	if inv.WarehouseID != 2 {
		t.Errorf("expected WarehouseID 2, got %d", inv.WarehouseID)
	}
	if inv.ProductID != 3 {
		t.Errorf("expected ProductID 3, got %d", inv.ProductID)
	}
	if inv.BasePrice != 100.0 {
		t.Errorf("expected BasePrice 100, got %f", inv.BasePrice)
	}
	if inv.SaleModel != inventory.SaleModelRetail {
		t.Errorf("expected SaleModel %q, got %q", inventory.SaleModelRetail, inv.SaleModel)
	}
	if inv.PromotionStatus != inventory.PromotionStatusPending {
		t.Errorf("expected PromotionStatus %q, got %q", inventory.PromotionStatusPending, inv.PromotionStatus)
	}
	if inv.MinOrderQty != 1 {
		t.Errorf("expected MinOrderQty 1, got %d", inv.MinOrderQty)
	}
	if inv.Condition != inventory.ConditionNew {
		t.Errorf("expected Condition %q, got %q", inventory.ConditionNew, inv.Condition)
	}
	if inv.VendorSaleStatus != inventory.VendorSaleStatusActive {
		t.Errorf("expected VendorSaleStatus %q, got %q", inventory.VendorSaleStatusActive, inv.VendorSaleStatus)
	}
	if inv.SystemSaleStatus != inventory.SystemSaleStatusActive {
		t.Errorf("expected SystemSaleStatus %q, got %q", inventory.SystemSaleStatusActive, inv.SystemSaleStatus)
	}
	if inv.PromotionID != nil {
		t.Errorf("expected PromotionID to be nil, got %v", inv.PromotionID)
	}
	if inv.FinalPrice != nil {
		t.Errorf("expected FinalPrice to be nil, got %v", inv.FinalPrice)
	}
	if inv.InstantQty != 0 {
		t.Errorf("expected InstantQty 0, got %d", inv.InstantQty)
	}
	if inv.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewInventory_ZeroBasePrice_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := inventory.NewInventory(1, 2, 3, 0)
	if err != inventory.ErrInvalidBasePrice {
		t.Errorf("expected ErrInvalidBasePrice, got %v", err)
	}
}

func TestNewInventory_NegativeBasePrice_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := inventory.NewInventory(1, 2, 3, -10.0)
	if err != inventory.ErrInvalidBasePrice {
		t.Errorf("expected ErrInvalidBasePrice, got %v", err)
	}
}

func TestApplyPromotion_Success(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	startAt := now
	endAt := now.Add(24 * time.Hour)
	promotionID := int64(42)
	finalPrice := 80.0

	err = inv.ApplyPromotion(promotionID, finalPrice, startAt, endAt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if inv.PromotionID == nil || *inv.PromotionID != promotionID {
		t.Errorf("expected PromotionID %d, got %v", promotionID, inv.PromotionID)
	}
	if inv.FinalPrice == nil || *inv.FinalPrice != finalPrice {
		t.Errorf("expected FinalPrice %f, got %v", finalPrice, inv.FinalPrice)
	}
	if inv.StartAt == nil || !inv.StartAt.Equal(startAt) {
		t.Errorf("expected StartAt %v, got %v", startAt, inv.StartAt)
	}
	if inv.EndAt == nil || !inv.EndAt.Equal(endAt) {
		t.Errorf("expected EndAt %v, got %v", endAt, inv.EndAt)
	}
	if inv.PromotionStatus != inventory.PromotionStatusPending {
		t.Errorf("expected PromotionStatus %q, got %q", inventory.PromotionStatusPending, inv.PromotionStatus)
	}
}

func TestApplyPromotion_AlreadyApplied_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = inv.ApplyPromotion(1, 80.0, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	err = inv.ApplyPromotion(2, 70.0, now, now.Add(24*time.Hour))
	if err != inventory.ErrPromotionAlreadyApplied {
		t.Errorf("expected ErrPromotionAlreadyApplied, got %v", err)
	}
}

func TestApplyPromotion_EndBeforeStart_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = inv.ApplyPromotion(1, 80.0, now.Add(24*time.Hour), now)
	if err != inventory.ErrInvalidPromotionDates {
		t.Errorf("expected ErrInvalidPromotionDates, got %v", err)
	}
}

func TestApplyPromotion_ZeroFinalPrice_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = inv.ApplyPromotion(1, 0, now, now.Add(24*time.Hour))
	if err != inventory.ErrInvalidFinalPrice {
		t.Errorf("expected ErrInvalidFinalPrice, got %v", err)
	}
}

func TestApplyPromotion_NegativeFinalPrice_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = inv.ApplyPromotion(1, -5.0, now, now.Add(24*time.Hour))
	if err != inventory.ErrInvalidFinalPrice {
		t.Errorf("expected ErrInvalidFinalPrice, got %v", err)
	}
}

func TestRemovePromotion_Success(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = inv.ApplyPromotion(1, 80.0, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	err = inv.RemovePromotion()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if inv.PromotionID != nil {
		t.Errorf("expected PromotionID to be nil, got %v", inv.PromotionID)
	}
	if inv.FinalPrice != nil {
		t.Errorf("expected FinalPrice to be nil, got %v", inv.FinalPrice)
	}
	if inv.StartAt != nil {
		t.Errorf("expected StartAt to be nil, got %v", inv.StartAt)
	}
	if inv.EndAt != nil {
		t.Errorf("expected EndAt to be nil, got %v", inv.EndAt)
	}
	if inv.PromotionStatus != inventory.PromotionStatusPending {
		t.Errorf("expected PromotionStatus %q, got %q", inventory.PromotionStatusPending, inv.PromotionStatus)
	}
}

func TestRemovePromotion_NoActivePromotion_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	err = inv.RemovePromotion()
	if err != inventory.ErrNoActivePromotion {
		t.Errorf("expected ErrNoActivePromotion, got %v", err)
	}
}

func TestUpdateInventory_Success(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	maxOrderQty := 100
	scheduledQty := map[string]int{"2026-07-06": 50}

	err = inv.UpdateInventory(10, scheduledQty, 5, &maxOrderQty)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if inv.InstantQty != 10 {
		t.Errorf("expected InstantQty 10, got %d", inv.InstantQty)
	}
	if len(inv.ScheduledQty) != 1 || inv.ScheduledQty["2026-07-06"] != 50 {
		t.Errorf("expected ScheduledQty map[2026-07-06:50], got %v", inv.ScheduledQty)
	}
	if inv.MinOrderQty != 5 {
		t.Errorf("expected MinOrderQty 5, got %d", inv.MinOrderQty)
	}
	if inv.MaxOrderQty == nil || *inv.MaxOrderQty != maxOrderQty {
		t.Errorf("expected MaxOrderQty %d, got %v", maxOrderQty, inv.MaxOrderQty)
	}
}

func TestUpdateInventory_NegativeInstantQty_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	maxOrderQty := 100
	err = inv.UpdateInventory(-1, nil, 5, &maxOrderQty)
	if err != inventory.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestUpdateInventory_ZeroMinOrderQty_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	maxOrderQty := 100
	err = inv.UpdateInventory(10, nil, 0, &maxOrderQty)
	if err != inventory.ErrInvalidMinOrderQty {
		t.Errorf("expected ErrInvalidMinOrderQty, got %v", err)
	}
}

func TestUpdateInventory_MaxOrderQtyLessThanMinOrderQty_ReturnsError(t *testing.T) {
	t.Parallel()

	inv, err := inventory.NewInventory(1, 2, 3, 100.0)
	if err != nil {
		t.Fatal(err)
	}

	maxOrderQty := 3
	err = inv.UpdateInventory(10, nil, 5, &maxOrderQty)
	if err != inventory.ErrInvalidMaxOrderQty {
		t.Errorf("expected ErrInvalidMaxOrderQty, got %v", err)
	}
}
