package promotion_test

import (
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

func TestListPromotions_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	listUC := promotionapp.NewListPromotionsUseCase(repo)

	for i := 0; i < 5; i++ {
		input := validCreateInput()
		input.Title = "Sale"
		createUC.Execute(input)
	}

	result, err := listUC.Execute(promotionapp.ListPromotionsInput{Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 5 {
		t.Errorf("expected Total 5, got %d", result.Total)
	}
	if len(result.Promotions) != 5 {
		t.Errorf("expected 5 promotions, got %d", len(result.Promotions))
	}
}

func TestListPromotions_FilterByStatus(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	listUC := promotionapp.NewListPromotionsUseCase(repo)

	for i := 0; i < 3; i++ {
		input := validCreateInput()
		input.Title = "Active"
		p, _ := createUC.Execute(input)
		p.Activate()
		repo.Save(p)
	}
	inactiveInput := validCreateInput()
	inactiveInput.Title = "Inactive"
	createUC.Execute(inactiveInput)

	active := promotion.PromotionStatusActive
	result, err := listUC.Execute(promotionapp.ListPromotionsInput{
		Status: &active, Page: 1, Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 3 {
		t.Errorf("expected 3 active, got %d", result.Total)
	}
}

func TestListPromotions_Pagination(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	listUC := promotionapp.NewListPromotionsUseCase(repo)

	for i := 0; i < 10; i++ {
		input := validCreateInput()
		input.Title = "Item"
		createUC.Execute(input)
	}

	result, err := listUC.Execute(promotionapp.ListPromotionsInput{Page: 1, Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 10 {
		t.Errorf("expected Total 10, got %d", result.Total)
	}
	if len(result.Promotions) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.Promotions))
	}
	if result.Page != 1 || result.Limit != 3 {
		t.Errorf("expected Page 1 Limit 3, got Page %d Limit %d", result.Page, result.Limit)
	}
}

func TestListPromotions_Search(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	createUC := promotionapp.NewCreatePromotionUseCase(repo)
	listUC := promotionapp.NewListPromotionsUseCase(repo)

	inputs := []string{"Summer Sale", "Winter Sale", "Clearance"}
	for _, title := range inputs {
		input := validCreateInput()
		input.Title = title
		createUC.Execute(input)
	}

	search := "Sale"
	result, err := listUC.Execute(promotionapp.ListPromotionsInput{
		Search: &search, Page: 1, Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 results for 'Sale', got %d", result.Total)
	}
}
