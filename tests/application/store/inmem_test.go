package store_test

import (
	"stock-service/internal/domain/store"
)

type inMemoryStoreRepository struct {
	stores map[int64]*store.Store
	lastID int64
}

func newInMemoryStoreRepository() *inMemoryStoreRepository {
	return &inMemoryStoreRepository{
		stores: make(map[int64]*store.Store),
	}
}

func (r *inMemoryStoreRepository) Save(s *store.Store) error {
	if s.ID == 0 {
		r.lastID++
		s.ID = r.lastID
	}
	r.stores[s.ID] = s
	return nil
}

func (r *inMemoryStoreRepository) FindByID(id int64) (*store.Store, error) {
	s, ok := r.stores[id]
	if !ok {
		return nil, store.ErrStoreNotFound
	}
	return s, nil
}

func (r *inMemoryStoreRepository) FindAll() ([]*store.Store, error) {
	result := make([]*store.Store, 0, len(r.stores))
	for _, s := range r.stores {
		result = append(result, s)
	}
	return result, nil
}

func (r *inMemoryStoreRepository) Delete(id int64) error {
	delete(r.stores, id)
	return nil
}
