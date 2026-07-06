package memory

import (
	"sync"

	categorydomain "stock-service/internal/domain/category"
)

type CategoryRepository struct {
	mu     sync.Mutex
	items  map[int64]*categorydomain.Category
	nextID int64
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{items: make(map[int64]*categorydomain.Category)}
}

func (r *CategoryRepository) Save(c *categorydomain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		r.nextID++
		c.ID = r.nextID
	}
	r.items[c.ID] = c
	return nil
}

func (r *CategoryRepository) FindByID(id int64) (*categorydomain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *CategoryRepository) FindByParentID(parentID int64) ([]*categorydomain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*categorydomain.Category
	for _, c := range r.items {
		if c.ParentID != nil && *c.ParentID == parentID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *CategoryRepository) FindAll() ([]*categorydomain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*categorydomain.Category, 0, len(r.items))
	for _, c := range r.items {
		result = append(result, c)
	}
	return result, nil
}

func (r *CategoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
