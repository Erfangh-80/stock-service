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

func (r *WarehouseRepository) FindAll(filter warehousedomain.WarehouseFilter) ([]*warehousedomain.Warehouse, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []*warehousedomain.Warehouse
	for _, w := range r.items {
		if filter.CreatedByUserID != nil && w.CreatedByUserID != *filter.CreatedByUserID {
			continue
		}
		if filter.IsPublic != nil && w.IsPublic != *filter.IsPublic {
			continue
		}
		result = append(result, w)
	}

	total := len(result)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= len(result) {
		return []*warehousedomain.Warehouse{}, total, nil
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], total, nil
}

func (r *WarehouseRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
