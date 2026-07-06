package memory

import (
	"sync"

	branddomain "stock-service/internal/domain/brand"
)

type BrandRepository struct {
	mu     sync.Mutex
	items  map[int64]*branddomain.Brand
	nextID int64
}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{items: make(map[int64]*branddomain.Brand)}
}

func (r *BrandRepository) Save(b *branddomain.Brand) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b.ID == 0 {
		r.nextID++
		b.ID = r.nextID
	}
	r.items[b.ID] = b
	return nil
}

func (r *BrandRepository) FindByID(id int64) (*branddomain.Brand, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *BrandRepository) FindAll() ([]*branddomain.Brand, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*branddomain.Brand, 0, len(r.items))
	for _, b := range r.items {
		result = append(result, b)
	}
	return result, nil
}

func (r *BrandRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
