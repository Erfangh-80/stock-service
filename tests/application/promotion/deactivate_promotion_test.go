package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestDeactivatePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewDeactivatePromotionUseCase(repo)

	p, _ := promotion.NewPromotion("Test")
	p.Activate()
	repo.Save(p)

	err := uc.Execute(p.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(p.ID)
	if saved.Status != promotion.PromotionStatusInactive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusInactive, saved.Status)
	}
}

func TestDeactivatePromotion_NotFound_ReturnsError(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewDeactivatePromotionUseCase(repo)

	err := uc.Execute(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
