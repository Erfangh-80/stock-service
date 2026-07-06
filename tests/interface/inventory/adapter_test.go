package inventoryinterface_test

import (
	"testing"
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	"stock-service/internal/domain/product"
	iface "stock-service/internal/interface"
	inventoryinterface "stock-service/internal/interface/inventory"
)

type mockCreate struct{ fn func(appinventory.CreateInventoryInput) (*inventory.Inventory, error) }
func (m *mockCreate) Execute(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockGet struct{ fn func(appinventory.GetInventoryInput) (*inventory.Inventory, error) }
func (m *mockGet) Execute(i appinventory.GetInventoryInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockList struct{ fn func(appinventory.ListInventoryInput) (*appinventory.ListInventoryOutput, error) }
func (m *mockList) Execute(i appinventory.ListInventoryInput) (*appinventory.ListInventoryOutput, error) { return m.fn(i) }

type mockDelete struct{ fn func(appinventory.DeleteInventoryInput) error }
func (m *mockDelete) Execute(i appinventory.DeleteInventoryInput) error { return m.fn(i) }

type mockSearch struct{ fn func(appinventory.SearchInventoryInput) (*appinventory.SearchInventoryOutput, error) }
func (m *mockSearch) Execute(i appinventory.SearchInventoryInput) (*appinventory.SearchInventoryOutput, error) { return m.fn(i) }

type mockApply struct{ fn func(appinventory.ApplyPromotionInput) (*inventory.Inventory, error) }
func (m *mockApply) Execute(i appinventory.ApplyPromotionInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockRemove struct{ fn func(appinventory.RemovePromotionInput) (*inventory.Inventory, error) }
func (m *mockRemove) Execute(i appinventory.RemovePromotionInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockUpdate struct{ fn func(appinventory.UpdateInventoryInput) (*inventory.Inventory, error) }
func (m *mockUpdate) Execute(i appinventory.UpdateInventoryInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockSuspendVendor struct{ fn func(appinventory.SuspendVendorSaleInput) (*inventory.Inventory, error) }
func (m *mockSuspendVendor) Execute(i appinventory.SuspendVendorSaleInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockCloseVendor struct{ fn func(appinventory.CloseVendorSaleInput) (*inventory.Inventory, error) }
func (m *mockCloseVendor) Execute(i appinventory.CloseVendorSaleInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockSuspendSystem struct{ fn func(appinventory.SuspendSystemSaleInput) (*inventory.Inventory, error) }
func (m *mockSuspendSystem) Execute(i appinventory.SuspendSystemSaleInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockCloseSystem struct{ fn func(appinventory.CloseSystemSaleInput) (*inventory.Inventory, error) }
func (m *mockCloseSystem) Execute(i appinventory.CloseSystemSaleInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockReserve struct{ fn func(appinventory.ReserveQuantityInput) (*inventory.Inventory, error) }
func (m *mockReserve) Execute(i appinventory.ReserveQuantityInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockRelease struct{ fn func(appinventory.ReleaseQuantityInput) (*inventory.Inventory, error) }
func (m *mockRelease) Execute(i appinventory.ReleaseQuantityInput) (*inventory.Inventory, error) { return m.fn(i) }

type mockLowStock struct{ fn func(appinventory.CheckLowStockInput) (*appinventory.CheckLowStockOutput, error) }
func (m *mockLowStock) Execute(i appinventory.CheckLowStockInput) (*appinventory.CheckLowStockOutput, error) { return m.fn(i) }

func nilAdapter() *inventoryinterface.Adapter {
	return inventoryinterface.NewAdapter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
}

func now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-07-05T12:00:00Z")
	return t
}

func baseInventory() *inventory.Inventory {
	return &inventory.Inventory{
		ID: 1, StoreID: 1, WarehouseID: 1, ProductID: 42,
		SaleModel: inventory.SaleModelRetail, BasePrice: 100,
		PromotionStatus: inventory.PromotionStatusPending,
		InstantQty: 0, MinOrderQty: 1, Condition: inventory.ConditionNew,
		VendorSaleStatus: inventory.VendorSaleStatusActive,
		SystemSaleStatus: inventory.SystemSaleStatusActive,
		CreatedAt: now(),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	inv := baseInventory()
	adapter := inventoryinterface.NewAdapter(
		&mockCreate{func(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) {
			if i.StoreID != 1 || i.BasePrice != 100 {
				t.Error("unexpected input")
			}
			inv.ID = 1
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.Create(inventoryinterface.CreateInventoryParams{
		StoreID: 1, WarehouseID: 1, ProductID: 42, BasePrice: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.BasePrice != 100 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestAdapter_Create_InvalidInput(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		&mockCreate{func(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrInvalidBasePrice
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.Create(inventoryinterface.CreateInventoryParams{
		StoreID: 1, WarehouseID: 1, ProductID: 42, BasePrice: 0,
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_Create_ProductNotFound(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		&mockCreate{func(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) {
			return nil, product.ErrProductNotFound
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.Create(inventoryinterface.CreateInventoryParams{
		StoreID: 1, WarehouseID: 1, ProductID: 999, BasePrice: 100,
	})
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_Get_Success(t *testing.T) {
	inv := baseInventory()
	adapter := inventoryinterface.NewAdapter(
		nil,
		&mockGet{func(i appinventory.GetInventoryInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestAdapter_Get_NotFound(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil,
		&mockGet{func(i appinventory.GetInventoryInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrInventoryNotFound
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_List_Success(t *testing.T) {
	inv := baseInventory()
	adapter := inventoryinterface.NewAdapter(
		nil, nil,
		&mockList{func(i appinventory.ListInventoryInput) (*appinventory.ListInventoryOutput, error) {
			return &appinventory.ListInventoryOutput{
				Items: []*inventory.Inventory{inv},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.List(inventoryinterface.ListInventoryParams{Page: 1, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestAdapter_Delete_Success(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil,
		&mockDelete{func(i appinventory.DeleteInventoryInput) error {
			if i.SaleID != 1 {
				t.Error("unexpected id")
			}
			return nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	err := adapter.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Delete_NotFound(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil,
		&mockDelete{func(i appinventory.DeleteInventoryInput) error {
			return inventory.ErrInventoryNotFound
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	err := adapter.Delete(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_Search_Success(t *testing.T) {
	inv := baseInventory()
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil,
		&mockSearch{func(i appinventory.SearchInventoryInput) (*appinventory.SearchInventoryOutput, error) {
			return &appinventory.SearchInventoryOutput{
				Items: []*inventory.Inventory{inv},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.Search(inventoryinterface.SearchInventoryParams{Query: "test", Page: 1, Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 result, got %d", resp.Total)
	}
}

func TestAdapter_ApplyPromotion_Success(t *testing.T) {
	inv := baseInventory()
	inv.PromotionID = ptr(int64(1))
	inv.FinalPrice = ptr(89.99)
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil,
		&mockApply{func(i appinventory.ApplyPromotionInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.ApplyPromotion(inventoryinterface.ApplyPromotionParams{
		SaleID: 1, PromotionID: 1, FinalPrice: 89.99, StartAt: now(), EndAt: now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.PromotionID == nil || *resp.PromotionID != 1 {
		t.Error("expected promotion ID 1")
	}
}

func TestAdapter_ApplyPromotion_Conflict(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil,
		&mockApply{func(i appinventory.ApplyPromotionInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrPromotionAlreadyApplied
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.ApplyPromotion(inventoryinterface.ApplyPromotionParams{
		SaleID: 1, PromotionID: 1, FinalPrice: 80, StartAt: now(), EndAt: now().Add(24 * time.Hour),
	})
	if err != iface.ErrConflict {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestAdapter_RemovePromotion_Success(t *testing.T) {
	inv := baseInventory()
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockRemove{func(i appinventory.RemovePromotionInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.RemovePromotion(inventoryinterface.RemovePromotionParams{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_UpdateInventory_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 10
	inv.MinOrderQty = 5
	maxQty := 100
	inv.MaxOrderQty = &maxQty
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil,
		&mockUpdate{func(i appinventory.UpdateInventoryInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.UpdateInventory(inventoryinterface.UpdateInventoryParams{
		SaleID: 1, InstantQty: 10, MinOrderQty: 5, MaxOrderQty: &maxQty,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.InstantQty != 10 || resp.MinOrderQty != 5 || *resp.MaxOrderQty != 100 {
		t.Error("unexpected response")
	}
}

func TestAdapter_SuspendVendorSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.VendorSaleStatus = inventory.VendorSaleStatusSuspended
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil,
		&mockSuspendVendor{func(i appinventory.SuspendVendorSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil,
	)

	resp, err := adapter.SuspendVendorSale(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.VendorSaleStatus != string(inventory.VendorSaleStatusSuspended) {
		t.Errorf("expected suspended, got %s", resp.VendorSaleStatus)
	}
}

func TestAdapter_CloseVendorSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.VendorSaleStatus = inventory.VendorSaleStatusClosed
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockCloseVendor{func(i appinventory.CloseVendorSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil,
	)

	resp, err := adapter.CloseVendorSale(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.VendorSaleStatus != string(inventory.VendorSaleStatusClosed) {
		t.Errorf("expected closed, got %s", resp.VendorSaleStatus)
	}
}

func TestAdapter_SuspendSystemSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.SystemSaleStatus = inventory.SystemSaleStatusSuspended
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockSuspendSystem{func(i appinventory.SuspendSystemSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil,
	)

	resp, err := adapter.SuspendSystemSale(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.SystemSaleStatus != string(inventory.SystemSaleStatusSuspended) {
		t.Errorf("expected suspended, got %s", resp.SystemSaleStatus)
	}
}

func TestAdapter_CloseSystemSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.SystemSaleStatus = inventory.SystemSaleStatusClosed
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockCloseSystem{func(i appinventory.CloseSystemSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil,
	)

	resp, err := adapter.CloseSystemSale(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.SystemSaleStatus != string(inventory.SystemSaleStatusClosed) {
		t.Errorf("expected closed, got %s", resp.SystemSaleStatus)
	}
}

func TestAdapter_ReserveQuantity_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 40
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockReserve{func(i appinventory.ReserveQuantityInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil,
	)

	resp, err := adapter.ReserveQuantity(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.InstantQty != 40 {
		t.Errorf("expected qty 40, got %d", resp.InstantQty)
	}
}

func TestAdapter_ReleaseQuantity_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 60
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockRelease{func(i appinventory.ReleaseQuantityInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil,
	)

	resp, err := adapter.ReleaseQuantity(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.InstantQty != 60 {
		t.Errorf("expected qty 60, got %d", resp.InstantQty)
	}
}

func TestAdapter_CheckLowStock_Success(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockLowStock{func(i appinventory.CheckLowStockInput) (*appinventory.CheckLowStockOutput, error) {
			return &appinventory.CheckLowStockOutput{IsLow: true, CurrentQty: 3}, nil
		}},
	)

	resp, err := adapter.CheckLowStock(1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsLow || resp.CurrentQty != 3 {
		t.Errorf("expected low=true, qty=3, got low=%v, qty=%d", resp.IsLow, resp.CurrentQty)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestAdapter_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil,
		&mockGet{func(i appinventory.GetInventoryInput) (*inventory.Inventory, error) {
			return nil, unknownErr{}
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	_, err := adapter.Get(1)
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func ptr[T any](v T) *T { return &v }
