package memory

import (
	"sync"

	productattributedomain "stock-service/internal/domain/product_attribute"
)

type ProductAttributeRepository struct {
	mu     sync.Mutex
	items  map[int64]*productattributedomain.ProductAttribute
	nextID int64
}

func NewProductAttributeRepository() *ProductAttributeRepository {
	return &ProductAttributeRepository{items: make(map[int64]*productattributedomain.ProductAttribute)}
}

func (r *ProductAttributeRepository) Save(pa *productattributedomain.ProductAttribute) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pa.ID == 0 {
		r.nextID++
		pa.ID = r.nextID
	}
	r.items[pa.ID] = pa
	return nil
}

func (r *ProductAttributeRepository) FindByID(id int64) (*productattributedomain.ProductAttribute, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ProductAttributeRepository) FindByProductID(productID int32) ([]*productattributedomain.ProductAttribute, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*productattributedomain.ProductAttribute
	for _, pa := range r.items {
		if pa.ProductID == productID {
			result = append(result, pa)
		}
	}
	return result, nil
}

func (r *ProductAttributeRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
