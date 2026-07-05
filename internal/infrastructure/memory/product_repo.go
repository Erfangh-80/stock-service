package memory

import (
	"strings"
	"sync"

	productdomain "stock-service/internal/domain/product"
)

type ProductRepository struct {
	mu    sync.Mutex
	items map[int32]*productdomain.Product
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{items: make(map[int32]*productdomain.Product)}
}

func (r *ProductRepository) Save(p *productdomain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[p.ID] = p
	return nil
}

func (r *ProductRepository) FindByID(id int32) (*productdomain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *ProductRepository) FindByTitle(query string) ([]*productdomain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*productdomain.Product
	for _, p := range r.items {
		if strings.Contains(p.TitleFa, query) || (p.TitleEn != nil && strings.Contains(*p.TitleEn, query)) {
			result = append(result, p)
		}
	}
	return result, nil
}

func SeedProducts(r *ProductRepository) {
	seeds := []struct {
		id      int32
		titleFa string
	}{
		{1, "محصول یک"},
		{2, "محصول دو"},
		{3, "محصول سه"},
		{10, "محصول ده"},
		{30, "محصول سی"},
		{42, "محصول چهل و دو"},
		{100, "محصول صد"},
	}
	for _, s := range seeds {
		p, err := productdomain.NewProduct(s.titleFa, 1, 1)
		if err != nil {
			panic(err)
		}
		p.ID = s.id
		r.items[p.ID] = p
	}
}
