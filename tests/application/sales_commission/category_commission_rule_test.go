package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

type inMemoryRuleRepo struct {
	rules  map[int64]*domainsalescommission.CategoryCommissionRule
	nextID int64
}

func newInMemoryRuleRepo() *inMemoryRuleRepo {
	return &inMemoryRuleRepo{
		rules:  make(map[int64]*domainsalescommission.CategoryCommissionRule),
		nextID: 1,
	}
}

func (r *inMemoryRuleRepo) Save(rule *domainsalescommission.CategoryCommissionRule) error {
	if rule.ID == 0 {
		rule.ID = r.nextID
		r.nextID++
	}
	r.rules[rule.ID] = rule
	return nil
}

func (r *inMemoryRuleRepo) FindByID(id int64) (*domainsalescommission.CategoryCommissionRule, error) {
	return r.rules[id], nil
}

func (r *inMemoryRuleRepo) FindAll(filter domainsalescommission.CategoryCommissionRuleFilter) ([]*domainsalescommission.CategoryCommissionRule, int, error) {
	var matched []*domainsalescommission.CategoryCommissionRule
	for _, rule := range r.rules {
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

func (r *inMemoryRuleRepo) Delete(id int64) error {
	delete(r.rules, id)
	return nil
}

func TestCreateCategoryCommissionRule_Success(t *testing.T) {
	repo := newInMemoryRuleRepo()
	uc := salescommission.NewCreateCategoryCommissionRuleUseCase(repo)

	rule, err := uc.Execute(salescommission.CreateCategoryCommissionRuleInput{
		CategoryID: 1, RatePercent: 10, MinPrice: 100,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rule.ID == 0 {
		t.Error("expected ID to be set")
	}
	if rule.CategoryID != 1 || rule.RatePercent != 10 || rule.MinPrice != 100 {
		t.Errorf("unexpected rule: %+v", rule)
	}
	if !rule.IsActive {
		t.Error("expected rule to be active by default")
	}
}

func TestGetCategoryCommissionRule_Success(t *testing.T) {
	repo := newInMemoryRuleRepo()
	rule, _ := domainsalescommission.NewCategoryCommissionRule(1, 10, 100)
	repo.Save(rule)

	uc := salescommission.NewGetCategoryCommissionRuleUseCase(repo)
	result, err := uc.Execute(salescommission.GetCategoryCommissionRuleInput{ID: rule.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.CategoryID != 1 {
		t.Errorf("expected CategoryID 1, got %d", result.CategoryID)
	}
}

func TestGetCategoryCommissionRule_NotFound(t *testing.T) {
	repo := newInMemoryRuleRepo()
	uc := salescommission.NewGetCategoryCommissionRuleUseCase(repo)

	_, err := uc.Execute(salescommission.GetCategoryCommissionRuleInput{ID: 999})
	if err != domainsalescommission.ErrRuleNotFound {
		t.Errorf("expected ErrRuleNotFound, got %v", err)
	}
}

func TestListCategoryCommissionRules(t *testing.T) {
	repo := newInMemoryRuleRepo()
	repo.Save(&domainsalescommission.CategoryCommissionRule{CategoryID: 1, RatePercent: 10, MinPrice: 100, IsActive: true})
	repo.Save(&domainsalescommission.CategoryCommissionRule{CategoryID: 2, RatePercent: 15, MinPrice: 200, IsActive: false})

	uc := salescommission.NewListCategoryCommissionRulesUseCase(repo)
	result, err := uc.Execute(salescommission.ListCategoryCommissionRulesInput{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(result.Rules))
	}
}

func TestUpdateCategoryCommissionRule_RatePercent(t *testing.T) {
	repo := newInMemoryRuleRepo()
	rule, _ := domainsalescommission.NewCategoryCommissionRule(1, 10, 100)
	repo.Save(rule)

	newRate := 15.0
	uc := salescommission.NewUpdateCategoryCommissionRuleUseCase(repo)
	result, err := uc.Execute(salescommission.UpdateCategoryCommissionRuleInput{
		ID: rule.ID, RatePercent: &newRate,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.RatePercent != 15 {
		t.Errorf("expected 15, got %f", result.RatePercent)
	}
}

func TestUpdateCategoryCommissionRule_Deactivate(t *testing.T) {
	repo := newInMemoryRuleRepo()
	rule, _ := domainsalescommission.NewCategoryCommissionRule(1, 10, 100)
	repo.Save(rule)

	active := false
	uc := salescommission.NewUpdateCategoryCommissionRuleUseCase(repo)
	result, err := uc.Execute(salescommission.UpdateCategoryCommissionRuleInput{
		ID: rule.ID, Activate: &active,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.IsActive {
		t.Error("expected rule to be inactive")
	}
}

func TestUpdateCategoryCommissionRule_NotFound(t *testing.T) {
	repo := newInMemoryRuleRepo()
	uc := salescommission.NewUpdateCategoryCommissionRuleUseCase(repo)

	_, err := uc.Execute(salescommission.UpdateCategoryCommissionRuleInput{ID: 999})
	if err != domainsalescommission.ErrRuleNotFound {
		t.Errorf("expected ErrRuleNotFound, got %v", err)
	}
}

func TestDeleteCategoryCommissionRule_Success(t *testing.T) {
	repo := newInMemoryRuleRepo()
	rule, _ := domainsalescommission.NewCategoryCommissionRule(1, 10, 100)
	repo.Save(rule)

	uc := salescommission.NewDeleteCategoryCommissionRuleUseCase(repo)
	err := uc.Execute(salescommission.DeleteCategoryCommissionRuleInput{ID: rule.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(rule.ID)
	if saved != nil {
		t.Error("expected rule to be deleted")
	}
}

func TestDeleteCategoryCommissionRule_NotFound(t *testing.T) {
	repo := newInMemoryRuleRepo()
	uc := salescommission.NewDeleteCategoryCommissionRuleUseCase(repo)
	err := uc.Execute(salescommission.DeleteCategoryCommissionRuleInput{ID: 999})
	if err != domainsalescommission.ErrRuleNotFound {
		t.Errorf("expected ErrRuleNotFound, got %v", err)
	}
}
