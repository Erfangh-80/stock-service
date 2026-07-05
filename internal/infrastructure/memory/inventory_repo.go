package memory

import (
	"sync"

	inventorydomain "stock-service/internal/domain/inventory"
)

type InventoryRepository struct {
	mu     sync.Mutex
	items  map[int64]*inventorydomain.Inventory
	nextID int64
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{items: make(map[int64]*inventorydomain.Inventory)}
}

func (r *InventoryRepository) Save(inv *inventorydomain.Inventory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if inv.ID == 0 {
		r.nextID++
		inv.ID = r.nextID
	}
	r.items[inv.ID] = inv
	return nil
}

func (r *InventoryRepository) FindByID(id int64) (*inventorydomain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *InventoryRepository) FindAll() ([]*inventorydomain.Inventory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*inventorydomain.Inventory, 0, len(r.items))
	for _, inv := range r.items {
		result = append(result, inv)
	}
	return result, nil
}

func (r *InventoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
