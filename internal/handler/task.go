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

type TaskHandler struct {
	svc  *service.TaskService
	Repo *repository.GormRepo
}

func NewTaskHandler(svc *service.TaskService, repo *repository.GormRepo) *TaskHandler {
	return &TaskHandler{svc: svc, Repo: repo}
}

// Register Chi routes
func (h *TaskHandler) Register(r chi.Router) {
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(&task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write audit log
	_ = h.Repo.WriteAudit("task.create", task.Title)

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	task, err := h.svc.Get(uint(id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.ID = uint(id)
	if err := h.svc.Update(&task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write audit log
	statusMsg := ""
	if task.Status != "" {
		statusMsg = fmt.Sprintf("status updated to %s", task.Status)
	}
	_ = h.Repo.WriteAudit("task.update", fmt.Sprintf("%s (%s)", task.Title, statusMsg))

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	task, _ := h.svc.Get(uint(id)) // optional: get task title before deletion
	if err := h.svc.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write audit log
	_ = h.Repo.WriteAudit("task.delete", task.Title)

	w.WriteHeader(http.StatusNoContent)
}
