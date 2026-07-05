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

func (r *SalesCommissionRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
