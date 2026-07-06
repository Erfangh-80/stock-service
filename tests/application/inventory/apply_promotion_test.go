package inventory_test

import (
	"testing"
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	promotiondomain "stock-service/internal/domain/promotion"
)

type mockPromotionRepo struct {
	promotions map[int64]*promotiondomain.Promotion
}

func newMockPromotionRepo() *mockPromotionRepo {
	return &mockPromotionRepo{promotions: make(map[int64]*promotiondomain.Promotion)}
}

func (r *mockPromotionRepo) Save(p *promotiondomain.Promotion) error {
	r.promotions[p.ID] = p
	return nil
}

func (r *mockPromotionRepo) FindByID(id int64) (*promotiondomain.Promotion, error) {
	p, ok := r.promotions[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *mockPromotionRepo) FindAll(filter promotiondomain.PromotionFilter) ([]*promotiondomain.Promotion, int, error) {
	var items []*promotiondomain.Promotion
	for _, p := range r.promotions {
		items = append(items, p)
	}
	return items, len(items), nil
}

func (r *mockPromotionRepo) FindByCouponCode(code string) (*promotiondomain.Promotion, error) {
	return nil, nil
}

func (r *mockPromotionRepo) Delete(id int64) error {
	delete(r.promotions, id)
	return nil
}

func activePromotion() *promotiondomain.Promotion {
	p, _ := promotiondomain.NewPromotion(promotiondomain.CreatePromotionInput{
		Title:         "Test Promo",
		DiscountType:  promotiondomain.DiscountTypePercentage,
		DiscountValue: 10,
	})
	p.ID = 100
	p.Activate()
	return p
}

func TestApplyPromotionUseCase_Success(t *testing.T) {
	repo := newInmemoryRepository()
	promoRepo := newMockPromotionRepo()

	promo := activePromotion()
	promoRepo.Save(promo)

	repo.Save(&inventory.Inventory{ID: 1, BasePrice: 1000})
	uc := appinventory.NewApplyPromotionUseCase(repo, promoRepo)

	now := time.Now()
	input := appinventory.ApplyPromotionInput{
		SaleID:      1,
		PromotionID: 100,
		FinalPrice:  900,
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
	if sale.FinalPrice == nil || *sale.FinalPrice != 900 {
		t.Fatalf("expected FinalPrice 900, got %v", sale.FinalPrice)
	}
}

func TestApplyPromotionUseCase_AutoCalculatePrice(t *testing.T) {
	repo := newInmemoryRepository()
	promoRepo := newMockPromotionRepo()

	promo := activePromotion()
	promoRepo.Save(promo)

	repo.Save(&inventory.Inventory{ID: 1, BasePrice: 2000})
	uc := appinventory.NewApplyPromotionUseCase(repo, promoRepo)

	now := time.Now()
	input := appinventory.ApplyPromotionInput{
		SaleID:      1,
		PromotionID: 100,
		FinalPrice:  0,
		StartAt:     now,
		EndAt:       now.Add(24 * time.Hour),
	}

	sale, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.FinalPrice == nil || *sale.FinalPrice != 1800 {
		t.Fatalf("expected auto-calculated FinalPrice 1800 (10%% off 2000), got %v", sale.FinalPrice)
	}
}

func TestApplyPromotionUseCase_PromotionNotFound(t *testing.T) {
	repo := newInmemoryRepository()
	promoRepo := newMockPromotionRepo()

	repo.Save(&inventory.Inventory{ID: 1})
	uc := appinventory.NewApplyPromotionUseCase(repo, promoRepo)

	input := appinventory.ApplyPromotionInput{
		SaleID:      1,
		PromotionID: 999,
		FinalPrice:  49.99,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(24 * time.Hour),
	}

	_, err := uc.Execute(input)
	if err != promotiondomain.ErrPromotionNotFound {
		t.Fatalf("expected ErrPromotionNotFound, got %v", err)
	}
}

func TestApplyPromotionUseCase_InventoryNotFound(t *testing.T) {
	repo := newInmemoryRepository()
	promoRepo := newMockPromotionRepo()

	promo := activePromotion()
	promoRepo.Save(promo)

	uc := appinventory.NewApplyPromotionUseCase(repo, promoRepo)

	input := appinventory.ApplyPromotionInput{
		SaleID:      999,
		PromotionID: 100,
		FinalPrice:  49.99,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(24 * time.Hour),
	}

	_, err := uc.Execute(input)
	if err != inventory.ErrInventoryNotFound {
		t.Fatalf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestApplyPromotionUseCase_AlreadyApplied(t *testing.T) {
	repo := newInmemoryRepository()
	promoRepo := newMockPromotionRepo()

	promo := activePromotion()
	promoRepo.Save(promo)

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
	uc := appinventory.NewApplyPromotionUseCase(repo, promoRepo)

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
