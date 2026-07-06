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
	mux.HandleFunc("PUT /api/v1/sales-commissions/{id}/max-price", h.UpdateMaxPrice)
	mux.HandleFunc("PUT /api/v1/sales-commissions/{id}/min-qty", h.UpdateMinQty)
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

	if err := h.adapter.UpdateMaxPrice(salescommissioninterface.UpdateMaxPriceParams{
		CommissionID: id,
		MaxPrice:     body.MaxPrice,
	}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	if err := h.adapter.UpdateMinQty(salescommissioninterface.UpdateMinQtyParams{
		CommissionID: id,
		MinQty:       body.MinQty,
	}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
