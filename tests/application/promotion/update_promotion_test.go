package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestUpdatePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	updateUC := promotionapp.NewUpdatePromotionUseCase(repo)

	created, err := createUC.Execute(validCreateInput())
	if err != nil {
		t.Fatal(err)
	}

	newTitle := "Updated Sale"
	updated, err := updateUC.Execute(promotionapp.UpdatePromotionInput{
		ID:    created.ID,
		Title: &newTitle,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Updated Sale" {
		t.Errorf("expected 'Updated Sale', got %q", updated.Title)
	}
}

func TestUpdatePromotion_NotFound(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewUpdatePromotionUseCase(repo)

	newTitle := "Nope"
	_, err := uc.Execute(promotionapp.UpdatePromotionInput{
		ID: 999, Title: &newTitle,
	})
	if err != promotion.ErrPromotionNotFound {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionNotFound, err)
	}
}

func TestUpdatePromotion_DuplicateCouponCode(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	updateUC := promotionapp.NewUpdatePromotionUseCase(repo)

	code1 := "CODE1"
	input1 := validCreateInput()
	input1.CouponCode = &code1
	p1, err := createUC.Execute(input1)
	if err != nil {
		t.Fatal(err)
	}

	code2 := "CODE2"
	input2 := validCreateInput()
	input2.Title = "Another"
	input2.CouponCode = &code2
	_, err = createUC.Execute(input2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = updateUC.Execute(promotionapp.UpdatePromotionInput{
		ID: p1.ID, CouponCode: &code2,
	})
	if err != promotion.ErrCouponCodeAlreadyExists {
		t.Errorf("expected %v, got %v", promotion.ErrCouponCodeAlreadyExists, err)
	}
}
