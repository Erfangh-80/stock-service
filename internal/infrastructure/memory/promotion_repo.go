package memory

import (
	"strings"
	"sync"

	promotiondomain "stock-service/internal/domain/promotion"
)

type PromotionRepository struct {
	mu     sync.Mutex
	items  map[int64]*promotiondomain.Promotion
	nextID int64
}

func NewPromotionRepository() *PromotionRepository {
	return &PromotionRepository{items: make(map[int64]*promotiondomain.Promotion)}
}

func (r *PromotionRepository) Save(p *promotiondomain.Promotion) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == 0 {
		r.nextID++
		p.ID = r.nextID
	}
	r.items[p.ID] = p
	return nil
}

func (r *PromotionRepository) FindByID(id int64) (*promotiondomain.Promotion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *PromotionRepository) FindAll(filter promotiondomain.PromotionFilter) ([]*promotiondomain.Promotion, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filtered []*promotiondomain.Promotion
	for _, p := range r.items {
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

func (r *PromotionRepository) FindByCouponCode(code string) (*promotiondomain.Promotion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.items {
		if p.CouponCode != nil && *p.CouponCode == code {
			return p, nil
		}
	}
	return nil, nil
}

func (r *PromotionRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
