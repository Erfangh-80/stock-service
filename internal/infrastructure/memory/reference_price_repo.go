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

func (r *ReferencePriceRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
