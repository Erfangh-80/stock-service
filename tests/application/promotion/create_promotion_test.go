package promotion_test

import (
	"strings"
	"testing"

	promotionapp "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
)

type inMemoryPromotionRepo struct {
	promotions map[int64]*promotion.Promotion
	nextID     int64
}

func newInMemoryPromotionRepo() *inMemoryPromotionRepo {
	return &inMemoryPromotionRepo{
		promotions: make(map[int64]*promotion.Promotion),
		nextID:     1,
	}
}

func (r *inMemoryPromotionRepo) Save(p *promotion.Promotion) error {
	if p.ID == 0 {
		p.ID = r.nextID
		r.nextID++
	}
	r.promotions[p.ID] = p
	return nil
}

func (r *inMemoryPromotionRepo) FindByID(id int64) (*promotion.Promotion, error) {
	p, ok := r.promotions[id]
	if !ok {
		return nil, promotion.ErrPromotionNotFound
	}
	return p, nil
}

func (r *inMemoryPromotionRepo) Delete(id int64) error {
	delete(r.promotions, id)
	return nil
}

func TestCreatePromotion_Success(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	p, err := uc.Execute("Summer Sale")
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
}

func TestCreatePromotion_EmptyTitle_ReturnsErrTitleRequired(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	_, err := uc.Execute("")
	if err != promotion.ErrTitleRequired {
		t.Errorf("expected %v, got %v", promotion.ErrTitleRequired, err)
	}
}

func TestCreatePromotion_TitleTooLong_ReturnsErrTitleTooLong(t *testing.T) {
	repo := newInMemoryPromotionRepo()
	uc := promotionapp.NewCreatePromotionUseCase(repo)

	longTitle := strings.Repeat("a", 256)
	_, err := uc.Execute(longTitle)
	if err != promotion.ErrTitleTooLong {
		t.Errorf("expected %v, got %v", promotion.ErrTitleTooLong, err)
	}
}
