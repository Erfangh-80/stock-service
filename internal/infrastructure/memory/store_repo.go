package memory

import (
	"sync"

	storedomain "stock-service/internal/domain/store"
)

type StoreRepository struct {
	mu     sync.Mutex
	stores map[int64]*storedomain.Store
	nextID int64
}

func NewStoreRepository() *StoreRepository {
	return &StoreRepository{stores: make(map[int64]*storedomain.Store)}
}

func (r *StoreRepository) Save(s *storedomain.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.ID == 0 {
		r.nextID++
		s.ID = r.nextID
	}
	r.stores[s.ID] = s
	return nil
}

func (r *StoreRepository) FindByID(id int64) (*storedomain.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stores[id], nil
}

func (r *StoreRepository) FindAll() ([]*storedomain.Store, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*storedomain.Store, 0, len(r.stores))
	for _, s := range r.stores {
		result = append(result, s)
	}
	return result, nil
}

func (r *StoreRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.stores, id)
	return nil
}
