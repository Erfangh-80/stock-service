package dto

import (
	"encoding/json"
	"net/http"

	iface "stock-service/internal/interface"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func EncodeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func HandleError(w http.ResponseWriter, err error) {
	switch err {
	case iface.ErrNotFound:
		EncodeJSON(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
	case iface.ErrInvalidInput:
		EncodeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid input"})
	case iface.ErrConflict:
		EncodeJSON(w, http.StatusConflict, ErrorResponse{Error: "conflict"})
	default:
		EncodeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
	}
}
