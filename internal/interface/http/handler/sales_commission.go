package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type SalesCommissionHandler struct {
	adapter *salescommissioninterface.Adapter
}

func NewSalesCommissionHandler(adapter *salescommissioninterface.Adapter) *SalesCommissionHandler {
	return &SalesCommissionHandler{adapter: adapter}
}

func (h *SalesCommissionHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/sales-commissions", h.Create)
	mux.HandleFunc("GET /api/v1/sales-commissions", h.List)
	mux.HandleFunc("GET /api/v1/sales-commissions/{id}", h.Get)
	mux.HandleFunc("DELETE /api/v1/sales-commissions/{id}", h.Delete)
	mux.HandleFunc("GET /api/v1/sales-commissions/by-inventory/{inventoryId}", h.GetByInventory)
	mux.HandleFunc("PUT /api/v1/sales-commissions/{id}/max-price", h.UpdateMaxPrice)
	mux.HandleFunc("PUT /api/v1/sales-commissions/{id}/min-qty", h.UpdateMinQty)
	mux.HandleFunc("POST /api/v1/sales-commissions/calculate", h.Calculate)
}

func (h *SalesCommissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateSalesCommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(salescommissioninterface.CreateSalesCommissionParams{
		InventoryID:              body.InventoryID,
		CategoryCommissionRuleID: body.CategoryCommissionRuleID,
		SaleModel:                body.SaleModel,
		RatePercent:              body.RatePercent,
		MinPrice:                 body.MinPrice,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *SalesCommissionHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *SalesCommissionHandler) GetByInventory(w http.ResponseWriter, r *http.Request) {
	inventoryID, err := strconv.ParseInt(r.PathValue("inventoryId"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid inventory_id"})
		return
	}

	result, err := h.adapter.GetByInventory(inventoryID)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *SalesCommissionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var inventoryID *int64
	if s := q.Get("inventory_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid inventory_id"})
			return
		}
		inventoryID = &v
	}

	var saleModel *string
	if s := q.Get("sale_model"); s != "" {
		saleModel = &s
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(salescommissioninterface.ListSalesCommissionsParams{
		InventoryID: inventoryID,
		SaleModel:   saleModel,
		Page:        page,
		Limit:       limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *SalesCommissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Delete(id); err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *SalesCommissionHandler) UpdateMaxPrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateMaxPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	if err := h.adapter.UpdateMaxPrice(struct {
		CommissionID int64
		MaxPrice     float64
	}{CommissionID: id, MaxPrice: body.MaxPrice}); err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *SalesCommissionHandler) UpdateMinQty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateMinQtyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	if err := h.adapter.UpdateMinQty(struct {
		CommissionID int64
		MinQty       int
	}{CommissionID: id, MinQty: body.MinQty}); err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *SalesCommissionHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var body dto.CalculateCommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Calculate(salescommissioninterface.CalculateCommissionParams{
		InventoryID: body.InventoryID,
		Quantity:    body.Quantity,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}
