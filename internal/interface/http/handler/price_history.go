package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	pricehistoryinterface "stock-service/internal/interface/price_history"
)

type PriceHistoryHandler struct {
	adapter *pricehistoryinterface.Adapter
}

func NewPriceHistoryHandler(adapter *pricehistoryinterface.Adapter) *PriceHistoryHandler {
	return &PriceHistoryHandler{adapter: adapter}
}

func (h *PriceHistoryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products/{productId}/price-history", h.Create)
	mux.HandleFunc("GET /api/v1/products/{productId}/price-history", h.List)
}

func (h *PriceHistoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	var body dto.CreatePriceHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(pricehistoryinterface.CreatePriceHistoryParams{
		ProductID: int32(productID), OldPrice: body.OldPrice,
		NewPrice: body.NewPrice, ChangedBy: body.ChangedBy,
		Description: body.Description,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *PriceHistoryHandler) List(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	result, err := h.adapter.List(int32(productID))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}
