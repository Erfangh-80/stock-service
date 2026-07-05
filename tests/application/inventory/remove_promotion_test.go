package inventory_test

import (
	"testing"
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestRemovePromotionUseCase_Success(t *testing.T) {
	repo := newInmemoryRepository()
	pid := int64(50)
	fp := 25.0
	now := time.Now()
	repo.Save(&inventory.Inventory{
		ID:              1,
		PromotionID:     &pid,
		FinalPrice:      &fp,
		StartAt:         &now,
		EndAt:           &now,
		PromotionStatus: inventory.PromotionStatusPending,
	})
	uc := appinventory.NewRemovePromotionUseCase(repo)

	input := appinventory.RemovePromotionInput{SaleID: 1}
	sale, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.PromotionID != nil {
		t.Fatal("expected PromotionID to be nil after removal")
	}
	if sale.FinalPrice != nil {
		t.Fatal("expected FinalPrice to be nil after removal")
	}
}

func TestRemovePromotionUseCase_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewRemovePromotionUseCase(repo)

	input := appinventory.RemovePromotionInput{SaleID: 999}
	_, err := uc.Execute(input)
	if err != inventory.ErrInventoryNotFound {
		t.Fatalf("expected errNotFound, got %v", err)
	}
}

func TestRemovePromotionUseCase_NoActivePromotion(t *testing.T) {
	repo := newInmemoryRepository()
	repo.Save(&inventory.Inventory{ID: 1})
	uc := appinventory.NewRemovePromotionUseCase(repo)

	input := appinventory.RemovePromotionInput{SaleID: 1}
	_, err := uc.Execute(input)
	if err != inventory.ErrNoActivePromotion {
		t.Fatalf("expected ErrNoActivePromotion, got %v", err)
	}
}
