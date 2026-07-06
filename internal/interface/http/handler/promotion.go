package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	promotioninterface "stock-service/internal/interface/promotion"
)

type PromotionHandler struct {
	adapter *promotioninterface.Adapter
}

func NewPromotionHandler(adapter *promotioninterface.Adapter) *PromotionHandler {
	return &PromotionHandler{adapter: adapter}
}

func (h *PromotionHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/promotions", h.Create)
	mux.HandleFunc("GET /api/v1/promotions", h.List)
	mux.HandleFunc("GET /api/v1/promotions/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/promotions/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/promotions/{id}", h.Delete)
	mux.HandleFunc("POST /api/v1/promotions/{id}/activate", h.Activate)
	mux.HandleFunc("POST /api/v1/promotions/{id}/deactivate", h.Deactivate)
}

// Create handles POST /api/v1/promotions — Admin-only
func (h *PromotionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body promotioninterface.CreatePromotionParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(body)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

// Get handles GET /api/v1/promotions/{id} — Admin-only
func (h *PromotionHandler) Get(w http.ResponseWriter, r *http.Request) {
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

// Update handles PUT /api/v1/promotions/{id} — Admin-only
func (h *PromotionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body promotioninterface.UpdatePromotionParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Update(id, body)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /api/v1/promotions/{id} — Admin-only
func (h *PromotionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// Activate handles POST /api/v1/promotions/{id}/activate — Salesperson-only
func (h *PromotionHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.adapter.Activate(id); err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "activated"})
}

// Deactivate handles POST /api/v1/promotions/{id}/deactivate — Salesperson-only
func (h *PromotionHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.adapter.Deactivate(id); err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "deactivated"})
}

// List handles GET /api/v1/promotions — Admin-only
func (h *PromotionHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	params := promotioninterface.ListPromotionsParams{
		Status:       strPtrOrNil(q.Get("status")),
		DiscountType: strPtrOrNil(q.Get("discount_type")),
		Search:       strPtrOrNil(q.Get("search")),
		Page:         page,
		Limit:        limit,
	}

	result, err := h.adapter.List(params)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

