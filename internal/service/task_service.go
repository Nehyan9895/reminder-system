package service

import (
	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
)

type TaskService struct {
	repo *repository.GormRepo
}

func NewTaskService(repo *repository.GormRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(task *models.Task) error {
	return s.repo.CreateTask(task)
}

func (s *TaskService) Get(id uint) (*models.Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s *TaskService) List() ([]models.Task, error) {
	return s.repo.ListTasks()
}

func (s *TaskService) Update(task *models.Task) error {
	return s.repo.UpdateTask(task)
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.DeleteTask(id)
}
