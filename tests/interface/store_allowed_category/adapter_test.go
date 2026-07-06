package storeallowedcategoryinterface_test

import (
	"testing"
	"time"

	app "stock-service/internal/application/store_allowed_category"
	domain "stock-service/internal/domain/store_allowed_category"
	iface "stock-service/internal/interface"
	storeallowedcategoryinterface "stock-service/internal/interface/store_allowed_category"
)

type mockCreateCategory struct {
	fn func(int64, int64) (*domain.StoreAllowedCategory, error)
}

func (m *mockCreateCategory) Execute(storeID, categoryID int64) (*domain.StoreAllowedCategory, error) {
	return m.fn(storeID, categoryID)
}

type mockGetCategory struct {
	fn func(app.GetStoreCategoryInput) (*domain.StoreAllowedCategory, error)
}

func (m *mockGetCategory) Execute(input app.GetStoreCategoryInput) (*domain.StoreAllowedCategory, error) {
	return m.fn(input)
}

type mockListCategories struct {
	fn func(app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error)
}

func (m *mockListCategories) Execute(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error) {
	return m.fn(input)
}

type mockApproveCategory struct {
	fn func(app.ApproveCategoryInput) error
}

func (m *mockApproveCategory) Execute(input app.ApproveCategoryInput) error {
	return m.fn(input)
}

type mockRejectCategory struct {
	fn func(app.RejectCategoryInput) error
}

func (m *mockRejectCategory) Execute(input app.RejectCategoryInput) error {
	return m.fn(input)
}

type mockDeleteCategory struct {
	fn func(app.DeleteStoreCategoryInput) error
}

func (m *mockDeleteCategory) Execute(input app.DeleteStoreCategoryInput) error {
	return m.fn(input)
}

type mockValidateCategory struct {
	fn func(app.ValidateCategoryExistsInput) error
}

func (m *mockValidateCategory) Execute(input app.ValidateCategoryExistsInput) error {
	return m.fn(input)
}

func baseCategory() *domain.StoreAllowedCategory {
	return &domain.StoreAllowedCategory{
		ID: 1, StoreID: 100, CategoryID: 200,
		Status:    domain.StatusPending,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func newTestAdapter(
	create *mockCreateCategory,
	get *mockGetCategory,
	list *mockListCategories,
	approve *mockApproveCategory,
	reject *mockRejectCategory,
	del *mockDeleteCategory,
	validate *mockValidateCategory,
) *storeallowedcategoryinterface.Adapter {
	if create == nil {
		create = &mockCreateCategory{}
	}
	if get == nil {
		get = &mockGetCategory{}
	}
	if list == nil {
		list = &mockListCategories{}
	}
	if approve == nil {
		approve = &mockApproveCategory{}
	}
	if reject == nil {
		reject = &mockRejectCategory{}
	}
	if del == nil {
		del = &mockDeleteCategory{}
	}
	if validate == nil {
		validate = &mockValidateCategory{}
	}
	return storeallowedcategoryinterface.NewAdapter(create, get, list, approve, reject, del, validate)
}

func TestAdapter_Create_Success(t *testing.T) {
	sac := baseCategory()
	adapter := newTestAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*domain.StoreAllowedCategory, error) {
			return sac, nil
		}},
		nil, nil, nil, nil, nil,
		&mockValidateCategory{func(app.ValidateCategoryExistsInput) error { return nil }},
	)
	resp, err := adapter.Create(storeallowedcategoryinterface.CreateCategoryParams{
		StoreID: 100, CategoryID: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.StoreID != 100 || resp.CategoryID != 200 {
		t.Error("unexpected response")
	}
	if resp.Status != "pending" {
		t.Errorf("expected Status 'pending', got %q", resp.Status)
	}
}

func TestAdapter_Get_Success(t *testing.T) {
	sac := baseCategory()
	adapter := newTestAdapter(
		nil,
		&mockGetCategory{func(input app.GetStoreCategoryInput) (*domain.StoreAllowedCategory, error) {
			return sac, nil
		}},
		nil, nil, nil, nil, nil,
	)
	resp, err := adapter.Get(storeallowedcategoryinterface.GetCategoryParams{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestAdapter_List_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil,
		&mockListCategories{func(input app.ListStoreCategoriesInput) (*app.ListStoreCategoriesOutput, error) {
			return &app.ListStoreCategoriesOutput{
				Categories: []*domain.StoreAllowedCategory{baseCategory()},
				Total:      1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil,
	)
	resp, err := adapter.List(storeallowedcategoryinterface.ListCategoriesParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Categories) != 1 || resp.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestAdapter_Approve_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil,
		&mockApproveCategory{func(input app.ApproveCategoryInput) error {
			return nil
		}},
		nil, nil, nil,
	)
	err := adapter.Approve(struct{ CategoryID int64 }{CategoryID: 1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Reject_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil,
		&mockRejectCategory{func(input app.RejectCategoryInput) error {
			if input.CategoryID != 1 {
				t.Error("unexpected input")
			}
			return nil
		}},
		nil, nil,
	)
	err := adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{CategoryID: 1, SupportNote: "bad"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Delete_Success(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil, nil,
		&mockDeleteCategory{func(input app.DeleteStoreCategoryInput) error {
			return nil
		}},
		nil,
	)
	err := adapter.Delete(storeallowedcategoryinterface.DeleteCategoryParams{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestAdapter_Create_ValidateCategoryNotFound_ReturnsNotFound(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockValidateCategory{func(input app.ValidateCategoryExistsInput) error {
			return domain.ErrCategoryNotFound
		}},
	)
	_, err := adapter.Create(storeallowedcategoryinterface.CreateCategoryParams{
		StoreID: 100, CategoryID: 999,
	})
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_Create_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := newTestAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*domain.StoreAllowedCategory, error) {
			return nil, unknownErr{}
		}},
		nil, nil, nil, nil, nil,
		&mockValidateCategory{func(app.ValidateCategoryExistsInput) error { return nil }},
	)
	_, err := adapter.Create(storeallowedcategoryinterface.CreateCategoryParams{
		StoreID: 100, CategoryID: 200,
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_Approve_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil,
		&mockApproveCategory{func(input app.ApproveCategoryInput) error {
			return unknownErr{}
		}},
		nil, nil, nil,
	)
	err := adapter.Approve(struct{ CategoryID int64 }{CategoryID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_Reject_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := newTestAdapter(
		nil, nil, nil, nil,
		&mockRejectCategory{func(input app.RejectCategoryInput) error {
			return unknownErr{}
		}},
		nil, nil,
	)
	err := adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{CategoryID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
