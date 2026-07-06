package inventoryhttp_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
	"stock-service/internal/interface/http/handler"
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

func TestHandler_Create_Success(t *testing.T) {
	inv := baseInventory()
	inv.ID = 1

	adapter := inventoryinterface.NewAdapter(
		&mockCreate{fn: func(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"store_id":1,"warehouse_id":1,"product_id":42,"base_price":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	raw, _ := io.ReadAll(resp.Body)
	var payload inventoryinterface.InventoryResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.ID != 1 || payload.BasePrice != 100 {
		t.Errorf("unexpected response: %+v", payload)
	}
}

func TestHandler_Create_InvalidInput(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		&mockCreate{fn: func(i appinventory.CreateInventoryInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrInvalidBasePrice
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"store_id":1,"warehouse_id":1,"product_id":42,"base_price":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	inv := baseInventory()

	adapter := inventoryinterface.NewAdapter(
		nil,
		&mockGet{fn: func(i appinventory.GetInventoryInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.ID != 1 {
		t.Errorf("expected ID 1, got %d", payload.ID)
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil,
		&mockGet{fn: func(i appinventory.GetInventoryInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrInventoryNotFound
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_InvalidID(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_List_Success(t *testing.T) {
	inv := baseInventory()

	adapter := inventoryinterface.NewAdapter(
		nil, nil,
		&mockList{fn: func(i appinventory.ListInventoryInput) (*appinventory.ListInventoryOutput, error) {
			return &appinventory.ListInventoryOutput{
				Items: []*inventory.Inventory{inv},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Total != 1 || len(payload.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(payload.Items))
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil,
		&mockDelete{fn: func(i appinventory.DeleteInventoryInput) error {
			return nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/inventory/1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil,
		&mockDelete{fn: func(i appinventory.DeleteInventoryInput) error {
			return inventory.ErrInventoryNotFound
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/inventory/999", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHandler_Search_Success(t *testing.T) {
	inv := baseInventory()

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil,
		&mockSearch{fn: func(i appinventory.SearchInventoryInput) (*appinventory.SearchInventoryOutput, error) {
			return &appinventory.SearchInventoryOutput{
				Items: []*inventory.Inventory{inv},
				Total: 1, Page: 1, Limit: 20,
			}, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/search?query=test", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Total != 1 {
		t.Errorf("expected 1 result, got %d", payload.Total)
	}
}

func TestHandler_Search_MissingQuery(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/search", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_ApplyPromotion_Success(t *testing.T) {
	inv := baseInventory()
	inv.PromotionID = ptr(int64(1))
	inv.FinalPrice = ptr(89.99)

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil,
		&mockApply{fn: func(i appinventory.ApplyPromotionInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"promotion_id":1,"final_price":89.99,"start_at":"2026-07-05T12:00:00Z","end_at":"2026-07-06T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/promotion", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.PromotionID == nil || *payload.PromotionID != 1 {
		t.Error("expected promotion ID 1")
	}
}

func TestHandler_ApplyPromotion_Conflict(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil,
		&mockApply{fn: func(i appinventory.ApplyPromotionInput) (*inventory.Inventory, error) {
			return nil, inventory.ErrPromotionAlreadyApplied
		}},
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"promotion_id":1,"final_price":80,"start_at":"2026-07-05T12:00:00Z","end_at":"2026-07-06T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/promotion", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", resp.StatusCode)
	}
}

func TestHandler_RemovePromotion_Success(t *testing.T) {
	inv := baseInventory()

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil,
		&mockRemove{fn: func(i appinventory.RemovePromotionInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/inventory/1/promotion", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_UpdateInventory_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 10
	inv.MinOrderQty = 5
	maxQty := 100
	inv.MaxOrderQty = &maxQty

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil,
		&mockUpdate{fn: func(i appinventory.UpdateInventoryInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"instant_qty":10,"min_order_qty":5,"max_order_qty":100}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/inventory/1/inventory", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.InstantQty != 10 || payload.MinOrderQty != 5 || payload.MaxOrderQty == nil || *payload.MaxOrderQty != 100 {
		t.Error("unexpected response")
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandler_SuspendVendorSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.VendorSaleStatus = inventory.VendorSaleStatusSuspended

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil,
		&mockSuspendVendor{fn: func(i appinventory.SuspendVendorSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/vendor/suspend", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.VendorSaleStatus != string(inventory.VendorSaleStatusSuspended) {
		t.Errorf("expected suspended, got %s", payload.VendorSaleStatus)
	}
}

func TestHandler_CloseVendorSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.VendorSaleStatus = inventory.VendorSaleStatusClosed

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockCloseVendor{fn: func(i appinventory.CloseVendorSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/vendor/close", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_SuspendSystemSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.SystemSaleStatus = inventory.SystemSaleStatusSuspended

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockSuspendSystem{fn: func(i appinventory.SuspendSystemSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/system/suspend", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_CloseSystemSale_Success(t *testing.T) {
	inv := baseInventory()
	inv.SystemSaleStatus = inventory.SystemSaleStatusClosed

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockCloseSystem{fn: func(i appinventory.CloseSystemSaleInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/system/close", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHandler_ReserveQuantity_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 40

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockReserve{fn: func(i appinventory.ReserveQuantityInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil, nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"quantity":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/reserve", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.InstantQty != 40 {
		t.Errorf("expected qty 40, got %d", payload.InstantQty)
	}
}

func TestHandler_ReleaseQuantity_Success(t *testing.T) {
	inv := baseInventory()
	inv.InstantQty = 60

	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockRelease{fn: func(i appinventory.ReleaseQuantityInput) (*inventory.Inventory, error) {
			return inv, nil
		}},
		nil,
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	body := `{"quantity":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventory/1/release", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.InstantQty != 60 {
		t.Errorf("expected qty 60, got %d", payload.InstantQty)
	}
}

func TestHandler_CheckLowStock_IsLow(t *testing.T) {
	adapter := inventoryinterface.NewAdapter(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&mockLowStock{fn: func(i appinventory.CheckLowStockInput) (*appinventory.CheckLowStockOutput, error) {
			return &appinventory.CheckLowStockOutput{IsLow: true, CurrentQty: 3}, nil
		}},
	)

	mux := http.NewServeMux()
	handler.NewInventoryHandler(adapter).Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory/1/low-stock?threshold=5", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var payload inventoryinterface.LowStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !payload.IsLow || payload.CurrentQty != 3 {
		t.Errorf("expected low=true, qty=3, got low=%v, qty=%d", payload.IsLow, payload.CurrentQty)
	}
}

func ptr[T any](v T) *T { return &v }
