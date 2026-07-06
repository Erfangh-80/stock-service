package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestDeletePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	deleteUC := promotionapp.NewDeletePromotionUseCase(repo)

	created, err := createUC.Execute(validCreateInput())
	if err != nil {
		t.Fatal(err)
	}

	err = deleteUC.Execute(promotionapp.DeletePromotionInput{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindByID(created.ID)
	if err != promotion.ErrPromotionNotFound {
		t.Errorf("expected ErrPromotionNotFound after delete, got %v", err)
	}
}

func TestDeletePromotion_NotFound(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewDeletePromotionUseCase(repo)

	err := uc.Execute(promotionapp.DeletePromotionInput{ID: 999})
	if err != promotion.ErrPromotionNotFound {
		t.Errorf("expected %v, got %v", promotion.ErrPromotionNotFound, err)
	}
}
