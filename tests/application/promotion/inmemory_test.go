package promotion_test

import (
	"strings"

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

func (r *inMemoryPromotionRepo) FindAll(filter promotion.PromotionFilter) ([]*promotion.Promotion, int, error) {
	var filtered []*promotion.Promotion
	for _, p := range r.promotions {
		if filter.Status != nil && p.Status != *filter.Status {
			continue
		}
		if filter.DiscountType != nil && p.DiscountType != *filter.DiscountType {
			continue
		}
		if filter.Search != nil {
			search := strings.ToLower(*filter.Search)
			if !strings.Contains(strings.ToLower(p.Title), search) {
				continue
			}
		}
		filtered = append(filtered, p)
	}
	total := len(filtered)
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	start := (page - 1) * limit
	if start > len(filtered) {
		return nil, total, nil
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (r *inMemoryPromotionRepo) FindByCouponCode(code string) (*promotion.Promotion, error) {
	for _, p := range r.promotions {
		if p.CouponCode != nil && *p.CouponCode == code {
			return p, nil
		}
	}
	return nil, nil
}

func (r *inMemoryPromotionRepo) Delete(id int64) error {
	delete(r.promotions, id)
	return nil
}
