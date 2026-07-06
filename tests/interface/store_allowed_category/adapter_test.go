package storeallowedcategoryinterface_test

import (
	"testing"
	"time"

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

type mockApproveCategory struct {
	fn func(int64) error
}

func (m *mockApproveCategory) Execute(categoryID int64) error {
	return m.fn(categoryID)
}

type mockRejectCategory struct {
	fn func(int64) error
}

func (m *mockRejectCategory) Execute(categoryID int64) error {
	return m.fn(categoryID)
}

func baseCategory() *domain.StoreAllowedCategory {
	return &domain.StoreAllowedCategory{
		ID: 1, StoreID: 100, CategoryID: 200,
		Status:    domain.StatusPending,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	sac := baseCategory()
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*domain.StoreAllowedCategory, error) {
			return sac, nil
		}},
		&mockApproveCategory{nil},
		&mockRejectCategory{nil},
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

func TestAdapter_Approve_Success(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{nil},
		&mockApproveCategory{func(categoryID int64) error {
			if categoryID != 1 {
				t.Error("unexpected input")
			}
			return nil
		}},
		&mockRejectCategory{nil},
	)
	err := adapter.Approve(storeallowedcategoryinterface.ApproveCategoryParams{CategoryID: 1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Reject_Success(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{nil},
		&mockApproveCategory{nil},
		&mockRejectCategory{func(categoryID int64) error {
			if categoryID != 1 {
				t.Error("unexpected input")
			}
			return nil
		}},
	)
	err := adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{CategoryID: 1})
	if err != nil {
		t.Fatal(err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestAdapter_Create_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{func(storeID, categoryID int64) (*domain.StoreAllowedCategory, error) {
			return nil, unknownErr{}
		}},
		&mockApproveCategory{nil},
		&mockRejectCategory{nil},
	)
	_, err := adapter.Create(storeallowedcategoryinterface.CreateCategoryParams{
		StoreID: 100, CategoryID: 200,
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_Approve_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{nil},
		&mockApproveCategory{func(categoryID int64) error {
			return unknownErr{}
		}},
		&mockRejectCategory{nil},
	)
	err := adapter.Approve(storeallowedcategoryinterface.ApproveCategoryParams{CategoryID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_Reject_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := storeallowedcategoryinterface.NewAdapter(
		&mockCreateCategory{nil},
		&mockApproveCategory{nil},
		&mockRejectCategory{func(categoryID int64) error {
			return unknownErr{}
		}},
	)
	err := adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{CategoryID: 1})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
