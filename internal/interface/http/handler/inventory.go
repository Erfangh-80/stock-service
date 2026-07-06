package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	inventoryinterface "stock-service/internal/interface/inventory"
)

type InventoryHandler struct {
	adapter *inventoryinterface.Adapter
}

func NewInventoryHandler(adapter *inventoryinterface.Adapter) *InventoryHandler {
	return &InventoryHandler{adapter: adapter}
}

func (h *InventoryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/inventory", h.Create)
	mux.HandleFunc("GET /api/v1/inventory", h.List)
	mux.HandleFunc("GET /api/v1/inventory/{id}", h.Get)
	mux.HandleFunc("DELETE /api/v1/inventory/{id}", h.Delete)
	mux.HandleFunc("GET /api/v1/inventory/search", h.Search)
	mux.HandleFunc("POST /api/v1/inventory/{id}/promotion", h.ApplyPromotion)
	mux.HandleFunc("DELETE /api/v1/inventory/{id}/promotion", h.RemovePromotion)
	mux.HandleFunc("PUT /api/v1/inventory/{id}/inventory", h.UpdateInventory)
	mux.HandleFunc("POST /api/v1/inventory/{id}/vendor/suspend", h.SuspendVendorSale)
	mux.HandleFunc("POST /api/v1/inventory/{id}/vendor/close", h.CloseVendorSale)
	mux.HandleFunc("POST /api/v1/inventory/{id}/system/suspend", h.SuspendSystemSale)
	mux.HandleFunc("POST /api/v1/inventory/{id}/system/close", h.CloseSystemSale)
	mux.HandleFunc("POST /api/v1/inventory/{id}/reserve", h.ReserveQuantity)
	mux.HandleFunc("POST /api/v1/inventory/{id}/release", h.ReleaseQuantity)
	mux.HandleFunc("GET /api/v1/inventory/{id}/low-stock", h.CheckLowStock)
}

func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(inventoryinterface.CreateInventoryParams{
		StoreID: body.StoreID, WarehouseID: body.WarehouseID,
		ProductID: body.ProductID, BasePrice: body.BasePrice,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *InventoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.Get(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var storeID *int64
	if s := q.Get("store_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid store_id"})
			return
		}
		storeID = &v
	}

	var productID *int32
	if s := q.Get("product_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product_id"})
			return
		}
		v32 := int32(v)
		productID = &v32
	}

	var vendorStatus *string
	if s := q.Get("vendor_sale_status"); s != "" {
		vendorStatus = &s
	}

	var systemStatus *string
	if s := q.Get("system_sale_status"); s != "" {
		systemStatus = &s
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(inventoryinterface.ListInventoryParams{
		StoreID: storeID, ProductID: productID,
		VendorSaleStatus: vendorStatus, SystemSaleStatus: systemStatus,
		Page: page, Limit: limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Delete(id); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	if query == "" {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "query is required"})
		return
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.Search(inventoryinterface.SearchInventoryParams{
		Query: query, Page: page, Limit: limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

// ApplyPromotion handles POST /api/v1/inventory/{id}/promotion — Salesperson-only
func (h *InventoryHandler) ApplyPromotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.ApplyPromotionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.ApplyPromotion(inventoryinterface.ApplyPromotionParams{
		SaleID: id, PromotionID: body.PromotionID,
		FinalPrice: body.FinalPrice, StartAt: body.StartAt, EndAt: body.EndAt,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

// RemovePromotion handles DELETE /api/v1/inventory/{id}/promotion — Salesperson-only
func (h *InventoryHandler) RemovePromotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.RemovePromotion(inventoryinterface.RemovePromotionParams{SaleID: id})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.UpdateInventory(inventoryinterface.UpdateInventoryParams{
		SaleID: id, InstantQty: body.InstantQty,
		ScheduledQty: body.ScheduledQty, MinOrderQty: body.MinOrderQty,
		MaxOrderQty: body.MaxOrderQty,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) SuspendVendorSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.SuspendVendorSale(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) CloseVendorSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.CloseVendorSale(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) SuspendSystemSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.SuspendSystemSale(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) CloseSystemSale(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.CloseSystemSale(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) ReserveQuantity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.ReserveQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.ReserveQuantity(id, body.Quantity)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) ReleaseQuantity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.ReleaseQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.ReleaseQuantity(id, body.Quantity)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) CheckLowStock(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	thresholdStr := r.URL.Query().Get("threshold")
	threshold := 10
	if thresholdStr != "" {
		threshold, err = strconv.Atoi(thresholdStr)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid threshold"})
			return
		}
	}

	result, err := h.adapter.CheckLowStock(id, threshold)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}
