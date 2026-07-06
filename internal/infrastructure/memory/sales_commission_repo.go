package memory

import (
	"sync"

	salescommissiondomain "stock-service/internal/domain/sales_commission"
)

type SalesCommissionRepository struct {
	mu     sync.Mutex
	items  map[int64]*salescommissiondomain.SalesCommission
	nextID int64
}

func NewSalesCommissionRepository() *SalesCommissionRepository {
	return &SalesCommissionRepository{items: make(map[int64]*salescommissiondomain.SalesCommission)}
}

func (r *SalesCommissionRepository) Save(sc *salescommissiondomain.SalesCommission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sc.ID == 0 {
		r.nextID++
		sc.ID = r.nextID
	}
	r.items[sc.ID] = sc
	return nil
}

func (r *SalesCommissionRepository) FindByID(id int64) (*salescommissiondomain.SalesCommission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *SalesCommissionRepository) FindByInventoryID(inventoryID int64) (*salescommissiondomain.SalesCommission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, sc := range r.items {
		if sc.InventoryID == inventoryID {
			return sc, nil
		}
	}
	return nil, nil
}

func (r *SalesCommissionRepository) FindAll(filter salescommissiondomain.SalesCommissionFilter) ([]*salescommissiondomain.SalesCommission, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var matched []*salescommissiondomain.SalesCommission
	for _, sc := range r.items {
		if filter.InventoryID != nil && sc.InventoryID != *filter.InventoryID {
			continue
		}
		if filter.SaleModel != nil && string(sc.SaleModel) != *filter.SaleModel {
			continue
		}
		matched = append(matched, sc)
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

func (r *SalesCommissionRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
