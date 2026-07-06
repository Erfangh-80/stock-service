package memory

import (
	"sync"

	salescommissiondomain "stock-service/internal/domain/sales_commission"
)

type CategoryCommissionRuleRepository struct {
	mu     sync.Mutex
	items  map[int64]*salescommissiondomain.CategoryCommissionRule
	nextID int64
}

func NewCategoryCommissionRuleRepository() *CategoryCommissionRuleRepository {
	return &CategoryCommissionRuleRepository{items: make(map[int64]*salescommissiondomain.CategoryCommissionRule)}
}

func (r *CategoryCommissionRuleRepository) Save(rule *salescommissiondomain.CategoryCommissionRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rule.ID == 0 {
		r.nextID++
		rule.ID = r.nextID
	}
	r.items[rule.ID] = rule
	return nil
}

func (r *CategoryCommissionRuleRepository) FindByID(id int64) (*salescommissiondomain.CategoryCommissionRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[id], nil
}

func (r *CategoryCommissionRuleRepository) FindAll(filter salescommissiondomain.CategoryCommissionRuleFilter) ([]*salescommissiondomain.CategoryCommissionRule, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var matched []*salescommissiondomain.CategoryCommissionRule
	for _, rule := range r.items {
		if filter.CategoryID != nil && rule.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.IsActive != nil && rule.IsActive != *filter.IsActive {
			continue
		}
		matched = append(matched, rule)
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

func (r *CategoryCommissionRuleRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}
