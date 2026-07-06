package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /api/v1/reference-prices", h.List)
	mux.HandleFunc("GET /api/v1/reference-prices/{id}", h.Get)
	mux.HandleFunc("DELETE /api/v1/reference-prices/{id}", h.Delete)
	mux.HandleFunc("GET /api/v1/reference-prices/by-product/{productId}", h.GetByProduct)
	mux.HandleFunc("GET /api/v1/reference-prices/by-product/{productId}/validate", h.Validate)
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

func (h *ReferencePriceHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *ReferencePriceHandler) GetByProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product_id"})
		return
	}

	result, err := h.adapter.GetByProduct(int32(productID))
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ReferencePriceHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

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

	var source *string
	if s := q.Get("source"); s != "" {
		source = &s
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(referencepriceinterface.ListReferencePricesParams{
		ProductID: productID,
		Source:    source,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ReferencePriceHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *ReferencePriceHandler) Validate(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product_id"})
		return
	}

	result, err := h.adapter.Validate(int32(productID))
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}
