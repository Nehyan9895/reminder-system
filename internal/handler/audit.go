package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
	"github.com/go-chi/chi/v5"
)

type AuditHandler struct {
	Repo *repository.GormRepo
}

func NewAuditHandler(r *repository.GormRepo) *AuditHandler {
	return &AuditHandler{Repo: r}
}

// Register all Audit endpoints
func (h *AuditHandler) Register(r chi.Router) {
	r.Get("/audit", h.List)
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	var logs []models.AuditLog
	if err := h.Repo.DB.Order("created_at desc").Limit(200).Find(&logs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(logs)
}
