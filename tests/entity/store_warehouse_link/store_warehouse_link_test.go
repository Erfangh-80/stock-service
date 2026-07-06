package storewarehouselink

import (
	"testing"

	"stock-service/internal/domain/store_warehouse_link"
)

func TestNewStoreWarehouseLink_SetsRelationTypePrimary(t *testing.T) {
	swl := storewarehouselink.NewStoreWarehouseLink(1, 2)
	if swl.StoreID != 1 {
		t.Errorf("expected StoreID %d, got %d", 1, swl.StoreID)
	}
	if swl.WarehouseID != 2 {
		t.Errorf("expected WarehouseID %d, got %d", 2, swl.WarehouseID)
	}
	if swl.RelationType != storewarehouselink.RelationTypePrimary {
		t.Errorf("expected RelationType %q, got %q", storewarehouselink.RelationTypePrimary, swl.RelationType)
	}
}

func TestChangeRelationType_UpdatesType(t *testing.T) {
	swl := storewarehouselink.NewStoreWarehouseLink(1, 2)
	err := swl.ChangeRelationType(storewarehouselink.RelationTypePrimary)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if swl.RelationType != storewarehouselink.RelationTypePrimary {
		t.Errorf("expected RelationType %q, got %q", storewarehouselink.RelationTypePrimary, swl.RelationType)
	}
}

func TestChangeRelationType_InvalidRelationType_ReturnsError(t *testing.T) {
	swl := storewarehouselink.NewStoreWarehouseLink(1, 2)
	err := swl.ChangeRelationType("invalid")
	if err != storewarehouselink.ErrInvalidRelationType {
		t.Errorf("expected ErrInvalidRelationType, got %v", err)
	}
}
