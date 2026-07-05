package inventory_test

import (
	"testing"
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestApplyPromotionUseCase_Success(t *testing.T) {
	repo := newInmemoryRepository()
	repo.Save(&inventory.Inventory{ID: 1})
	uc := appinventory.NewApplyPromotionUseCase(repo)

	now := time.Now()
	input := appinventory.ApplyPromotionInput{
		SaleID:      1,
		PromotionID: 100,
		FinalPrice:  49.99,
		StartAt:     now,
		EndAt:       now.Add(24 * time.Hour),
	}

	sale, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.PromotionID == nil || *sale.PromotionID != 100 {
		t.Fatalf("expected PromotionID 100, got %v", sale.PromotionID)
	}
	if sale.FinalPrice == nil || *sale.FinalPrice != 49.99 {
		t.Fatalf("expected FinalPrice 49.99, got %v", sale.FinalPrice)
	}
}

func TestApplyPromotionUseCase_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewApplyPromotionUseCase(repo)

	input := appinventory.ApplyPromotionInput{
		SaleID:      999,
		PromotionID: 100,
		FinalPrice:  49.99,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(24 * time.Hour),
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrInventoryNotFound {
		t.Fatalf("expected errNotFound, got %v", err)
	}
}

func TestApplyPromotionUseCase_AlreadyApplied(t *testing.T) {
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
	uc := appinventory.NewApplyPromotionUseCase(repo)

	input := appinventory.ApplyPromotionInput{
		SaleID:      1,
		PromotionID: 100,
		FinalPrice:  49.99,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(24 * time.Hour),
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrPromotionAlreadyApplied {
		t.Fatalf("expected ErrPromotionAlreadyApplied, got %v", err)
	}
}
