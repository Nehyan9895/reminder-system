package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
	"github.com/Nehyan9895/reminder-system/internal/service"
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

	// --- Validation for uniqueness ---
	var rules []models.ReminderRule
	if err := h.Repo.DB.Find(&rules).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, rr := range rules {
		switch in.RuleType {
		case "at_due":
			if rr.RuleType == "at_due" {
				http.Error(w, "only one at_due rule is allowed", http.StatusBadRequest)
				return
			}
		case "before_due":
			var newParams service.BeforeDueParams
			if err := json.Unmarshal([]byte(in.Params), &newParams); err != nil {
				http.Error(w, "invalid params JSON", http.StatusBadRequest)
				return
			}

			if rr.RuleType == "before_due" {
				var existingParams service.BeforeDueParams
				if err := json.Unmarshal([]byte(rr.Params), &existingParams); err == nil {
					if existingParams.MinutesBefore == newParams.MinutesBefore {
						http.Error(w, fmt.Sprintf("before_due rule with %d minutes already exists", newParams.MinutesBefore), http.StatusBadRequest)
						return
					}
				}
			}
		case "interval":
			var newParams service.IntervalParams
			if err := json.Unmarshal([]byte(in.Params), &newParams); err != nil {
				http.Error(w, "invalid params JSON", http.StatusBadRequest)
				return
			}

			if rr.RuleType == "interval" {
				var existingParams service.IntervalParams
				if err := json.Unmarshal([]byte(rr.Params), &existingParams); err == nil {
					if existingParams.IntervalMin == newParams.IntervalMin {
						http.Error(w, fmt.Sprintf("interval rule with %d minutes already exists", newParams.IntervalMin), http.StatusBadRequest)
						return
					}
				}
			}
		}
	}

	// --- Persist ---
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

	// --- Validation for uniqueness ---
	var rules []models.ReminderRule
	if err := h.Repo.DB.Find(&rules).Error; err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, existing := range rules {
		// skip self
		if existing.ID == rr.ID {
			continue
		}

		switch in.RuleType {
		case "at_due":
			if existing.RuleType == "at_due" {
				http.Error(w, "only one at_due rule is allowed", http.StatusBadRequest)
				return
			}
		case "before_due":
			var newParams service.BeforeDueParams
			if err := json.Unmarshal([]byte(in.Params), &newParams); err != nil {
				http.Error(w, "invalid params JSON", http.StatusBadRequest)
				return
			}

			if existing.RuleType == "before_due" {
				var existingParams service.BeforeDueParams
				if err := json.Unmarshal([]byte(existing.Params), &existingParams); err == nil {
					if existingParams.MinutesBefore == newParams.MinutesBefore {
						http.Error(w, fmt.Sprintf("before_due rule with %d minutes already exists", newParams.MinutesBefore), http.StatusBadRequest)
						return
					}
				}
			}
		case "interval":
			var newParams service.IntervalParams
			if err := json.Unmarshal([]byte(in.Params), &newParams); err != nil {
				http.Error(w, "invalid params JSON", http.StatusBadRequest)
				return
			}

			if existing.RuleType == "interval" {
				var existingParams service.IntervalParams
				if err := json.Unmarshal([]byte(existing.Params), &existingParams); err == nil {
					if existingParams.IntervalMin == newParams.IntervalMin {
						http.Error(w, fmt.Sprintf("interval rule with %d minutes already exists", newParams.IntervalMin), http.StatusBadRequest)
						return
					}
				}
			}
		}
	}

	// --- Update values ---
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
