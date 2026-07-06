package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestGetPromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	getUC := promotionapp.NewGetPromotionUseCase(repo)

	created, err := createUC.Execute(validCreateInput())
	if err != nil {
		t.Fatal(err)
	}

	got, err := getUC.Execute(promotionapp.GetPromotionInput{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
}

func TestGetPromotion_NotFound(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	getUC := promotionapp.NewGetPromotionUseCase(repo)

	_, err := getUC.Execute(promotionapp.GetPromotionInput{ID: 999})
	if err != promotion.ErrPromotionNotFound {
		t.Errorf("expected ErrPromotionNotFound, got %v", err)
	}
}
