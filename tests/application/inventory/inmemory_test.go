package inventory_test

import (
	"sync"

	"stock-service/internal/domain/inventory"
)

type inmemoryRepository struct {
	mu     sync.Mutex
	store  map[int64]*inventory.Inventory
	nextID int64
}

func newInmemoryRepository() *inmemoryRepository {
	return &inmemoryRepository{
		store:  make(map[int64]*inventory.Inventory),
		nextID: 1,
	}
}

func (r *inmemoryRepository) Save(sale *inventory.Inventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sale.ID == 0 {
		sale.ID = r.nextID
		r.nextID++
	}
	r.store[sale.ID] = sale
	return nil
}

func (r *inmemoryRepository) FindByID(id int64) (*inventory.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	sale, ok := r.store[id]
	if !ok {
		return nil, inventory.ErrInventoryNotFound
	}
	return sale, nil
}

func (r *inmemoryRepository) FindAll() ([]*inventory.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*inventory.Inventory, 0, len(r.store))
	for _, inv := range r.store {
		result = append(result, inv)
	}
	return result, nil
}

func (r *inmemoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}
