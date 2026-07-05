package memory

import (
	"sync"

	warehousedomain "stock-service/internal/domain/warehouse"
)

type WarehouseRepository struct {
	mu     sync.Mutex
	items  map[int64]*warehousedomain.Warehouse
	nextID int64
}

func NewWarehouseRepository() *WarehouseRepository {
	return &WarehouseRepository{items: make(map[int64]*warehousedomain.Warehouse)}
}

func (r *WarehouseRepository) Save(w *warehousedomain.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == 0 {
		r.nextID++
		w.ID = r.nextID
	}
	r.items[w.ID] = w
	return nil
}

func (r *WarehouseRepository) FindByID(id int64) (*warehousedomain.Warehouse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *WarehouseRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
