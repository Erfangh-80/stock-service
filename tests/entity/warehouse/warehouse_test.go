package warehouse

import (
	"strings"
	"testing"

	"stock-service/internal/domain/warehouse"
)

func TestNewWarehouse_ValidInputs_Succeeds(t *testing.T) {
	w, err := warehouse.NewWarehouse(1, "Main Warehouse")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.CreatedByUserID != 1 {
		t.Errorf("expected CreatedByUserID %d, got %d", 1, w.CreatedByUserID)
	}
	if w.WarehouseName != "Main Warehouse" {
		t.Errorf("expected WarehouseName %q, got %q", "Main Warehouse", w.WarehouseName)
	}
}

func TestNewWarehouse_EmptyName_ReturnsErrWarehouseNameRequired(t *testing.T) {
	_, err := warehouse.NewWarehouse(1, "")
	if err != warehouse.ErrWarehouseNameRequired {
		t.Errorf("expected %v, got %v", warehouse.ErrWarehouseNameRequired, err)
	}
}

func TestNewWarehouse_NameTooLong_ReturnsErrWarehouseNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 256)
	_, err := warehouse.NewWarehouse(1, longName)
	if err != warehouse.ErrWarehouseNameTooLong {
		t.Errorf("expected %v, got %v", warehouse.ErrWarehouseNameTooLong, err)
	}
}

func TestMakePublic_SetsIsPublicTrue(t *testing.T) {
	w, _ := warehouse.NewWarehouse(1, "Test")
	w.MakePublic()
	if !w.IsPublic {
		t.Error("expected IsPublic to be true")
	}
}

func TestMakePrivate_SetsIsPublicFalse(t *testing.T) {
	w, _ := warehouse.NewWarehouse(1, "Test")
	w.MakePublic()
	w.MakePrivate()
	if w.IsPublic {
		t.Error("expected IsPublic to be false")
	}
}

func TestUpdatePhone_SetsPhone(t *testing.T) {
	w, _ := warehouse.NewWarehouse(1, "Test")
	phone := "1234567890"
	w.UpdatePhone(&phone)
	if w.Phone == nil {
		t.Fatal("expected Phone to be set, got nil")
	}
	if *w.Phone != "1234567890" {
		t.Errorf("expected Phone %q, got %q", "1234567890", *w.Phone)
	}
}

func TestUpdateContactPhone_SetsContactPhone(t *testing.T) {
	w, _ := warehouse.NewWarehouse(1, "Test")
	phone := "0987654321"
	w.UpdateContactPhone(&phone)
	if w.ContactPhone == nil {
		t.Fatal("expected ContactPhone to be set, got nil")
	}
	if *w.ContactPhone != "0987654321" {
		t.Errorf("expected ContactPhone %q, got %q", "0987654321", *w.ContactPhone)
	}
}

func TestUpdateCollectionMethod_SetsCollectionMethod(t *testing.T) {
	w, _ := warehouse.NewWarehouse(1, "Test")
	w.UpdateCollectionMethod("pickup")
	if w.CollectionMethod != "pickup" {
		t.Errorf("expected CollectionMethod %q, got %q", "pickup", w.CollectionMethod)
	}
}
