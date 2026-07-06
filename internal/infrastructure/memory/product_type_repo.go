package memory

import (
	"sort"
	"sync"

	producttypedomain "stock-service/internal/domain/product_type"
)

type ProductTypeRepository struct {
	mu     sync.Mutex
	items  map[int64]*producttypedomain.ProductType
	nextID int64
}

func NewProductTypeRepository() *ProductTypeRepository {
	return &ProductTypeRepository{items: make(map[int64]*producttypedomain.ProductType)}
}

func (r *ProductTypeRepository) Save(pt *producttypedomain.ProductType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pt.ID == 0 {
		r.nextID++
		pt.ID = r.nextID
	}
	r.items[pt.ID] = pt
	return nil
}

func (r *ProductTypeRepository) FindByID(id int64) (*producttypedomain.ProductType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ProductTypeRepository) FindByProductID(productID int32) ([]*producttypedomain.ProductType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*producttypedomain.ProductType
	for _, pt := range r.items {
		if pt.ProductID == productID {
			result = append(result, pt)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SortOrder < result[j].SortOrder
	})
	return result, nil
}

func (r *ProductTypeRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
