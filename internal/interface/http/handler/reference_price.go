package handler

import (
	"encoding/json"
	"net/http"

	"stock-service/internal/interface/http/dto"
	referencepriceinterface "stock-service/internal/interface/reference_price"
)

type ReferencePriceHandler struct {
	adapter *referencepriceinterface.Adapter
}

func NewReferencePriceHandler(adapter *referencepriceinterface.Adapter) *ReferencePriceHandler {
	return &ReferencePriceHandler{adapter: adapter}
}

func (h *ReferencePriceHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/reference-prices", h.Create)
}

func (h *ReferencePriceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateReferencePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(referencepriceinterface.CreateReferencePriceParams{
		ProductID: body.ProductID,
		Price:     body.Price,
		Source:    body.Source,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}
