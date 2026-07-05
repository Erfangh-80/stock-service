package memory

import (
	"sync"

	storeallowedcategorydomain "stock-service/internal/domain/store_allowed_category"
)

type StoreCategoryRepository struct {
	mu     sync.Mutex
	items  map[int64]*storeallowedcategorydomain.StoreAllowedCategory
	nextID int64
}

func NewStoreCategoryRepository() *StoreCategoryRepository {
	return &StoreCategoryRepository{items: make(map[int64]*storeallowedcategorydomain.StoreAllowedCategory)}
}

func (r *StoreCategoryRepository) Save(sac *storeallowedcategorydomain.StoreAllowedCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sac.ID == 0 {
		r.nextID++
		sac.ID = r.nextID
	}
	r.items[sac.ID] = sac
	return nil
}

func (r *StoreCategoryRepository) FindByID(id int64) (*storeallowedcategorydomain.StoreAllowedCategory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *StoreCategoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
