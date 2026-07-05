package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestActivatePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewActivatePromotionUseCase(repo)

	p, _ := promotion.NewPromotion("Test")
	repo.Save(p)

	err := uc.Execute(p.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(p.ID)
	if saved.Status != promotion.PromotionStatusActive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusActive, saved.Status)
	}
}

func TestActivatePromotion_NotFound_ReturnsError(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewActivatePromotionUseCase(repo)

	err := uc.Execute(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
