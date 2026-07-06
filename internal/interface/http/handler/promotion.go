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
	mux.HandleFunc("GET /api/v1/promotions/{id}", h.Get)
	mux.HandleFunc("POST /api/v1/promotions/{id}/activate", h.Activate)
	mux.HandleFunc("POST /api/v1/promotions/{id}/deactivate", h.Deactivate)
}

func (h *PromotionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreatePromotionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(promotioninterface.CreatePromotionParams{Title: body.Title})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

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
