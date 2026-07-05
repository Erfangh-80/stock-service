package memory

import (
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

func (r *PromotionRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
