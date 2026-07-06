package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func validCreateInput() promotionapp.CreatePromotionInput {
	return promotionapp.CreatePromotionInput{
		Title:         "Summer Sale",
		DiscountType:  promotion.DiscountTypePercentage,
		DiscountValue: 10,
	}
}

func TestCreatePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	p, err := uc.Execute(validCreateInput())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID == 0 {
		t.Error("expected ID to be set")
	}
	if p.Title != "Summer Sale" {
		t.Errorf("expected Title %q, got %q", "Summer Sale", p.Title)
	}
	if p.Status != promotion.PromotionStatusInactive {
		t.Errorf("expected Status %q, got %q", promotion.PromotionStatusInactive, p.Status)
	}
	if p.DiscountType != promotion.DiscountTypePercentage {
		t.Errorf("expected DiscountType %q, got %q", promotion.DiscountTypePercentage, p.DiscountType)
	}
}

func TestCreatePromotion_EmptyTitle_ReturnsErrTitleRequired(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	input := validCreateInput()
	input.Title = ""
	_, err := uc.Execute(input)
	if err != promotion.ErrTitleRequired {
		t.Errorf("expected %v, got %v", promotion.ErrTitleRequired, err)
	}
}

func TestCreatePromotion_InvalidDiscountType_ReturnsErr(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	input := validCreateInput()
	input.DiscountType = "invalid"
	_, err := uc.Execute(input)
	if err != promotion.ErrInvalidDiscountType {
		t.Errorf("expected %v, got %v", promotion.ErrInvalidDiscountType, err)
	}
}

func TestCreatePromotion_CouponCode_AlreadyExists(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	code := "SUMMER10"
	input1 := validCreateInput()
	input1.CouponCode = &code
	_, err := uc.Execute(input1)
	if err != nil {
		t.Fatal(err)
	}

	input2 := validCreateInput()
	input2.Title = "Another Sale"
	input2.CouponCode = &code
	_, err = uc.Execute(input2)
	if err != promotion.ErrCouponCodeAlreadyExists {
		t.Errorf("expected %v, got %v", promotion.ErrCouponCodeAlreadyExists, err)
	}
}
