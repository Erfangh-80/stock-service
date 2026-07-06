package memory

import (
	"strings"
	"sync"

	productdomain "stock-service/internal/domain/product"
)

type ProductRepository struct {
	mu    sync.Mutex
	items map[int32]*productdomain.Product
	nextID int32
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{items: make(map[int32]*productdomain.Product)}
}

func (r *ProductRepository) Save(p *productdomain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == 0 {
		r.nextID++
		p.ID = r.nextID
	}
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

func (r *ProductRepository) FindAll(filter productdomain.ProductFilter) ([]*productdomain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []*productdomain.Product
	for _, p := range r.items {
		if filter.OwnerType != nil && p.OwnerType != *filter.OwnerType {
			continue
		}
		if filter.OwnerID != nil && (p.OwnerID == nil || *p.OwnerID != *filter.OwnerID) {
			continue
		}
		if filter.Status != nil && p.Status != *filter.Status {
			continue
		}
		if filter.CategoryID != nil && p.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.BrandID != nil && p.BrandID != *filter.BrandID {
			continue
		}
		if filter.Search != nil {
			q := *filter.Search
			if !strings.Contains(p.TitleFa, q) && (p.TitleEn == nil || !strings.Contains(*p.TitleEn, q)) {
				continue
			}
		}
		result = append(result, p)
	}

	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	start := (page - 1) * limit
	if start >= len(result) {
		return nil, nil
	}

	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

func (r *ProductRepository) Count(filter productdomain.ProductFilter) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := 0
	for _, p := range r.items {
		if filter.OwnerType != nil && p.OwnerType != *filter.OwnerType {
			continue
		}
		if filter.OwnerID != nil && (p.OwnerID == nil || *p.OwnerID != *filter.OwnerID) {
			continue
		}
		if filter.Status != nil && p.Status != *filter.Status {
			continue
		}
		if filter.CategoryID != nil && p.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.BrandID != nil && p.BrandID != *filter.BrandID {
			continue
		}
		if filter.Search != nil {
			q := *filter.Search
			if !strings.Contains(p.TitleFa, q) && (p.TitleEn == nil || !strings.Contains(*p.TitleEn, q)) {
				continue
			}
		}
		count++
	}
	return count, nil
}

func SeedProducts(r *ProductRepository) {
	seeds := []struct {
		id      int32
		titleFa string
		slug    string
	}{
		{1, "محصول یک", "product-1"},
		{2, "محصول دو", "product-2"},
		{3, "محصول سه", "product-3"},
		{10, "محصول ده", "product-10"},
		{30, "محصول سی", "product-30"},
		{42, "محصول چهل و دو", "product-42"},
		{100, "محصول صد", "product-100"},
	}
	for _, s := range seeds {
		p, err := productdomain.NewProduct(s.titleFa, 1, 1, productdomain.WithSlug(s.slug))
		if err != nil {
			panic(err)
		}
		p.ID = s.id
		r.items[p.ID] = p
	}
}
