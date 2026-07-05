package storeallowedcategory

import (
	"testing"

	"stock-service/internal/domain/store_allowed_category"
)

func TestNewStoreAllowedCategory_SetsStatusPending(t *testing.T) {
	sac := storeallowedcategory.NewStoreAllowedCategory(1, 2)
	if sac.StoreID != 1 {
		t.Errorf("expected StoreID %d, got %d", 1, sac.StoreID)
	}
	if sac.CategoryID != 2 {
		t.Errorf("expected CategoryID %d, got %d", 2, sac.CategoryID)
	}
	if sac.Status != storeallowedcategory.StatusPending {
		t.Errorf("expected Status %q, got %q", storeallowedcategory.StatusPending, sac.Status)
	}
}

func TestApprove_SetsStatusApproved(t *testing.T) {
	sac := storeallowedcategory.NewStoreAllowedCategory(1, 2)
	sac.Approve()
	if sac.Status != storeallowedcategory.StatusApproved {
		t.Errorf("expected Status %q, got %q", storeallowedcategory.StatusApproved, sac.Status)
	}
}

func TestReject_SetsStatusRejected(t *testing.T) {
	sac := storeallowedcategory.NewStoreAllowedCategory(1, 2)
	sac.Reject()
	if sac.Status != storeallowedcategory.StatusRejected {
		t.Errorf("expected Status %q, got %q", storeallowedcategory.StatusRejected, sac.Status)
	}
}
