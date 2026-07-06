package promotion

import (
	"strings"
	"testing"
	"time"

	"stock-service/internal/domain/promotion"
)

func validInput() promotion.CreatePromotionInput {
	return promotion.CreatePromotionInput{
		Title:         "Summer Sale",
		DiscountType:  promotion.DiscountTypePercentage,
		DiscountValue: 10,
	}
}

func TestNewPromotion_Valid_SetsStatusInactive(t *testing.T) {
	p, err := promotion.NewPromotion(validInput())
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
	input := validInput()
	input.Title = ""
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrTitleRequired {
		t.Errorf("expected %v, got %v", promotion.ErrTitleRequired, err)
	}
}

func TestNewPromotion_TitleTooLong_ReturnsErrTitleTooLong(t *testing.T) {
	input := validInput()
	input.Title = strings.Repeat("a", 256)
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrTitleTooLong {
		t.Errorf("expected %v, got %v", promotion.ErrTitleTooLong, err)
	}
}

func TestNewPromotion_EmptyDiscountType_ReturnsErr(t *testing.T) {
	input := validInput()
	input.DiscountType = ""
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrDiscountTypeRequired {
		t.Errorf("expected %v, got %v", promotion.ErrDiscountTypeRequired, err)
	}
}

func TestNewPromotion_InvalidDiscountType_ReturnsErr(t *testing.T) {
	input := validInput()
	input.DiscountType = "invalid"
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrInvalidDiscountType {
		t.Errorf("expected %v, got %v", promotion.ErrInvalidDiscountType, err)
	}
}

func TestNewPromotion_ZeroDiscountValue_ReturnsErr(t *testing.T) {
	input := validInput()
	input.DiscountValue = 0
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrDiscountValueRequired {
		t.Errorf("expected %v, got %v", promotion.ErrDiscountValueRequired, err)
	}
}

func TestNewPromotion_PercentageTooHigh_ReturnsErr(t *testing.T) {
	input := validInput()
	input.DiscountValue = 150
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrDiscountValueTooHigh {
		t.Errorf("expected %v, got %v", promotion.ErrDiscountValueTooHigh, err)
	}
}

func TestNewPromotion_InvalidDates_ReturnsErr(t *testing.T) {
	input := validInput()
	start := mustTime("2026-07-06T12:00:00Z")
	end := mustTime("2026-07-05T12:00:00Z")
	input.StartAt = &start
	input.EndAt = &end
	_, err := promotion.NewPromotion(input)
	if err != promotion.ErrInvalidPromotionDates {
		t.Errorf("expected %v, got %v", promotion.ErrInvalidPromotionDates, err)
	}
}

func TestActivate_SetsStatusActive(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.Activate()
	if p.Status != promotion.PromotionStatusActive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusActive, p.Status)
	}
}

func TestDeactivate_SetsStatusInactive(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.Activate()
	p.Deactivate()
	if p.Status != promotion.PromotionStatusInactive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusInactive, p.Status)
	}
}

func TestIsActive(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	if p.IsActive() {
		t.Error("expected inactive")
	}
	p.Activate()
	if !p.IsActive() {
		t.Error("expected active")
	}
}

func TestIsExpired(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	if p.IsExpired() {
		t.Error("expected not expired")
	}
	past := mustTime("2020-01-01T00:00:00Z")
	p.EndAt = &past
	if !p.IsExpired() {
		t.Error("expected expired")
	}
}

func TestCanApply_Inactive_ReturnsErr(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	err := p.CanApply()
	if err != promotion.ErrPromotionNotActive {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionNotActive, err)
	}
}

func TestCanApply_Expired_ReturnsErr(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.Activate()
	past := mustTime("2020-01-01T00:00:00Z")
	p.EndAt = &past
	err := p.CanApply()
	if err != promotion.ErrPromotionExpired {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionExpired, err)
	}
}

func TestCanApply_UsageLimit_ReturnsErr(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	limit := 1
	p.UsageLimit = &limit
	p.UsedCount = 1
	p.Activate()
	err := p.CanApply()
	if err != promotion.ErrPromotionUsageLimitExceeded {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionUsageLimitExceeded, err)
	}
}

func TestCanApply_BudgetExceeded_ReturnsErr(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	budget := 1000.0
	p.Budget = &budget
	p.BudgetSpent = 1000
	p.Activate()
	err := p.CanApply()
	if err != promotion.ErrPromotionBudgetExceeded {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionBudgetExceeded, err)
	}
}

func TestCanApply_Success(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.Activate()
	err := p.CanApply()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCalculateDiscountPrice_Percentage(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	price := p.CalculateDiscountPrice(1000)
	if price != 900 {
		t.Errorf("expected 900, got %f", price)
	}
}

func TestCalculateDiscountPrice_PercentageWithCap(t *testing.T) {
	input := validInput()
	input.DiscountValue = 50
	cap := 200.0
	input.MaxDiscountAmount = &cap
	p, _ := promotion.NewPromotion(input)
	price := p.CalculateDiscountPrice(1000)
	if price != 800 {
		t.Errorf("expected 800, got %f", price)
	}
}

func TestCalculateDiscountPrice_FixedAmount(t *testing.T) {
	input := validInput()
	input.DiscountType = promotion.DiscountTypeFixedAmount
	input.DiscountValue = 300
	p, _ := promotion.NewPromotion(input)
	price := p.CalculateDiscountPrice(1000)
	if price != 700 {
		t.Errorf("expected 700, got %f", price)
	}
}

func TestCalculateDiscountPrice_MinPurchase(t *testing.T) {
	input := validInput()
	input.MinPurchase = float64Ptr(500)
	p, _ := promotion.NewPromotion(input)
	price := p.CalculateDiscountPrice(200)
	if price != 200 {
		t.Errorf("expected 200 (no discount below min), got %f", price)
	}
}

func TestIsEligibleForStore(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.EligibleStoreIDs = []int64{1, 2, 3}
	if p.IsEligibleForStore(1) != true {
		t.Error("expected eligible")
	}
	if p.IsEligibleForStore(4) != false {
		t.Error("expected not eligible")
	}
}

func TestIsEligibleForStore_EmptyList(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	if p.IsEligibleForStore(99) != true {
		t.Error("expected eligible when no restrictions")
	}
}

func TestUpdate_Partial(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	newTitle := "Updated Sale"
	err := p.Update(promotion.UpdatePromotionInput{Title: &newTitle})
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Updated Sale" {
		t.Errorf("expected 'Updated Sale', got %q", p.Title)
	}
	if p.DiscountType != promotion.DiscountTypePercentage {
		t.Errorf("discount type should remain unchanged")
	}
}

func TestRecordUsage(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	p.RecordUsage()
	if p.UsedCount != 1 {
		t.Errorf("expected UsedCount 1, got %d", p.UsedCount)
	}
}

func TestSpendBudget(t *testing.T) {
	p, _ := promotion.NewPromotion(validInput())
	budget := 1000.0
	p.Budget = &budget
	p.SpendBudget(500)
	if p.BudgetSpent != 500 {
		t.Errorf("expected BudgetSpent 500, got %f", p.BudgetSpent)
	}
}

func mustTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func float64Ptr(f float64) *float64 {
	return &f
}
