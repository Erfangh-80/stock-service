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

func (r *StoreCategoryRepository) FindAll(filter storeallowedcategorydomain.StoreCategoryFilter) ([]*storeallowedcategorydomain.StoreAllowedCategory, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var matched []*storeallowedcategorydomain.StoreAllowedCategory
	for _, sac := range r.items {
		if filter.StoreID != nil && sac.StoreID != *filter.StoreID {
			continue
		}
		matched = append(matched, sac)
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

func (r *StoreCategoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
