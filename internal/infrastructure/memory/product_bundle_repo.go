package memory

import (
	"sync"

	productbundledomain "stock-service/internal/domain/product_bundle"
)

type ProductBundleRepository struct {
	mu     sync.Mutex
	items  map[int64]*productbundledomain.ProductBundle
	nextID int64
}

func NewProductBundleRepository() *ProductBundleRepository {
	return &ProductBundleRepository{items: make(map[int64]*productbundledomain.ProductBundle)}
}

func (r *ProductBundleRepository) Save(pb *productbundledomain.ProductBundle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pb.ID == 0 {
		r.nextID++
		pb.ID = r.nextID
	}
	r.items[pb.ID] = pb
	return nil
}

func (r *ProductBundleRepository) FindByID(id int64) (*productbundledomain.ProductBundle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ProductBundleRepository) FindByProductID(productID int32) ([]*productbundledomain.ProductBundle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*productbundledomain.ProductBundle
	for _, pb := range r.items {
		if pb.ProductID == productID {
			result = append(result, pb)
		}
	}
	return result, nil
}

func (r *ProductBundleRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
