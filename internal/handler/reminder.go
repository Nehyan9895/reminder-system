package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ReminderHandler struct {
	Repo *repository.GormRepo
}

func NewReminderHandler(r *repository.GormRepo) *ReminderHandler {
	return &ReminderHandler{Repo: r}
}

// Register all Reminder endpoints
func (h *ReminderHandler) Register(r chi.Router) {
	r.Route("/rules", func(r chi.Router) {
		r.Post("/", h.CreateRule)
		r.Get("/", h.ListRules)
		r.Get("/{id}", h.GetRule)
		r.Put("/{id}", h.UpdateRule)
		r.Delete("/{id}", h.DeleteRule)
		r.Post("/{id}/activate", h.Activate)
		r.Post("/{id}/deactivate", h.Deactivate)
	})
}

func (h *ReminderHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var in models.ReminderRule
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Repo.CreateRule(&in); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.Repo.WriteAudit("rule.create", in.Name)
	json.NewEncoder(w).Encode(in)
}

func (h *ReminderHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	rules, _ := h.Repo.ListRules()
	json.NewEncoder(w).Encode(rules)
}

func (h *ReminderHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	rr, err := h.Repo.GetRuleByID(uint(id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(rr)
}

func (h *ReminderHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var in models.ReminderRule
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rr, err := h.Repo.GetRuleByID(uint(id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	rr.Name = in.Name
	rr.Params = in.Params
	rr.RuleType = in.RuleType
	if err := h.Repo.UpdateRule(rr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = h.Repo.WriteAudit("rule.update", rr.Name)
	json.NewEncoder(w).Encode(rr)
}

func (h *ReminderHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = h.Repo.DeleteRule(uint(id))
	_ = h.Repo.WriteAudit("rule.delete", strconv.Itoa(id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ReminderHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = h.Repo.SetRuleActive(uint(id), true)
	_ = h.Repo.WriteAudit("rule.activate", strconv.Itoa(id))
	w.WriteHeader(http.StatusNoContent)
}

func (h *ReminderHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = h.Repo.SetRuleActive(uint(id), false)
	_ = h.Repo.WriteAudit("rule.deactivate", strconv.Itoa(id))
	w.WriteHeader(http.StatusNoContent)
}
