package memory

import (
	"sync"

	referencepricedomain "stock-service/internal/domain/reference_price"
)

type ReferencePriceRepository struct {
	mu     sync.Mutex
	items  map[int64]*referencepricedomain.ReferencePrice
	nextID int64
}

func NewReferencePriceRepository() *ReferencePriceRepository {
	return &ReferencePriceRepository{items: make(map[int64]*referencepricedomain.ReferencePrice)}
}

func (r *ReferencePriceRepository) Save(rp *referencepricedomain.ReferencePrice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rp.ID == 0 {
		r.nextID++
		rp.ID = r.nextID
	}
	r.items[rp.ID] = rp
	return nil
}

func (r *ReferencePriceRepository) FindByID(id int64) (*referencepricedomain.ReferencePrice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ReferencePriceRepository) FindByProductID(productID int32) (*referencepricedomain.ReferencePrice, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rp := range r.items {
		if rp.ProductID == productID {
			return rp, nil
		}
	}
	return nil, nil
}

func (r *ReferencePriceRepository) FindAll(filter referencepricedomain.ReferencePriceFilter) ([]*referencepricedomain.ReferencePrice, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var matched []*referencepricedomain.ReferencePrice
	for _, rp := range r.items {
		if filter.ProductID != nil && rp.ProductID != *filter.ProductID {
			continue
		}
		if filter.Source != nil && rp.Source != *filter.Source {
			continue
		}
		matched = append(matched, rp)
	}

	total := len(matched)

	page, limit := filter.Page, filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= len(matched) {
		return nil, total, nil
	}
	end := start + limit
	if end > len(matched) {
		end = len(matched)
	}

	return matched[start:end], total, nil
}

func (r *ReferencePriceRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
