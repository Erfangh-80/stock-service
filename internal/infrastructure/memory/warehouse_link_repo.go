package memory

import (
	"sync"

	storewarehouselinkdomain "stock-service/internal/domain/store_warehouse_link"
)

type WarehouseLinkRepository struct {
	mu     sync.Mutex
	items  map[int64]*storewarehouselinkdomain.StoreWarehouseLink
	nextID int64
}

func NewWarehouseLinkRepository() *WarehouseLinkRepository {
	return &WarehouseLinkRepository{items: make(map[int64]*storewarehouselinkdomain.StoreWarehouseLink)}
}

func (r *WarehouseLinkRepository) Save(swl *storewarehouselinkdomain.StoreWarehouseLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if swl.ID == 0 {
		r.nextID++
		swl.ID = r.nextID
	}
	r.items[swl.ID] = swl
	return nil
}

func (r *WarehouseLinkRepository) FindByID(id int64) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *WarehouseLinkRepository) FindAll(filter storewarehouselinkdomain.WarehouseLinkFilter) ([]*storewarehouselinkdomain.StoreWarehouseLink, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var matched []*storewarehouselinkdomain.StoreWarehouseLink
	for _, swl := range r.items {
		if filter.StoreID != nil && swl.StoreID != *filter.StoreID {
			continue
		}
		if filter.WarehouseID != nil && swl.WarehouseID != *filter.WarehouseID {
			continue
		}
		matched = append(matched, swl)
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

func (r *WarehouseLinkRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
