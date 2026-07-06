package memory

import (
	"sort"
	"sync"

	productimagedomain "stock-service/internal/domain/product_image"
)

type ProductImageRepository struct {
	mu     sync.Mutex
	items  map[int64]*productimagedomain.ProductImage
	nextID int64
}

func NewProductImageRepository() *ProductImageRepository {
	return &ProductImageRepository{items: make(map[int64]*productimagedomain.ProductImage)}
}

func (r *ProductImageRepository) Save(pi *productimagedomain.ProductImage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pi.ID == 0 {
		r.nextID++
		pi.ID = r.nextID
	}
	r.items[pi.ID] = pi
	return nil
}

func (r *ProductImageRepository) FindByID(id int64) (*productimagedomain.ProductImage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ProductImageRepository) FindByProductID(productID int32) ([]*productimagedomain.ProductImage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*productimagedomain.ProductImage
	for _, pi := range r.items {
		if pi.ProductID == productID {
			result = append(result, pi)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SortOrder < result[j].SortOrder
	})
	return result, nil
}

func (r *ProductImageRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
