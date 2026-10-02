package lead

import (
	"encoding/json"
	lead "lead-automation/internal/lead/model"

	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input lead.Lead

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateLead(r.Context(), &input)
	if err != nil {
		http.Error(w, "failed to create lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(created)
}