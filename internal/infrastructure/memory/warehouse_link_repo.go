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

func (r *WarehouseLinkRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
