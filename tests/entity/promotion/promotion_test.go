package promotion

import (
	"strings"
	"testing"

	"stock-service/internal/domain/promotion"
)

func TestNewPromotion_ValidTitle_SetsStatusInactive(t *testing.T) {
	p, err := promotion.NewPromotion("Summer Sale")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Status != promotion.PromotionStatusInactive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusInactive, p.Status)
	}
	if p.Title != "Summer Sale" {
		t.Errorf("expected Title %q, got %q", "Summer Sale", p.Title)
	}
}

func TestNewPromotion_EmptyTitle_ReturnsErrTitleRequired(t *testing.T) {
	_, err := promotion.NewPromotion("")
	if err != promotion.ErrTitleRequired {
		t.Errorf("expected %v, got %v", promotion.ErrTitleRequired, err)
	}
}

func TestNewPromotion_TitleTooLong_ReturnsErrTitleTooLong(t *testing.T) {
	longTitle := strings.Repeat("a", 256)
	_, err := promotion.NewPromotion(longTitle)
	if err != promotion.ErrTitleTooLong {
		t.Errorf("expected %v, got %v", promotion.ErrTitleTooLong, err)
	}
}

func TestActivate_SetsStatusActive(t *testing.T) {
	p, _ := promotion.NewPromotion("Test")
	p.Activate()
	if p.Status != promotion.PromotionStatusActive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusActive, p.Status)
	}
}

func TestDeactivate_SetsStatusInactive(t *testing.T) {
	p, _ := promotion.NewPromotion("Test")
	p.Activate()
	p.Deactivate()
	if p.Status != promotion.PromotionStatusInactive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusInactive, p.Status)
	}
}
