package repository

import (
	"github.com/Nehyan9895/reminder-system/internal/models"
)

func (r *GormRepo) ListPendingTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Where("status = ?", "pending").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *GormRepo) ListTasks() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *GormRepo) GetTaskByID(id uint) (*models.Task, error) {
	var t models.Task
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepo) CreateTask(t *models.Task) error {
	return r.DB.Create(t).Error
}

func (r *GormRepo) UpdateTask(t *models.Task) error {
	return r.DB.Save(t).Error
}

func (r *GormRepo) DeleteTask(id uint) error {
	return r.DB.Delete(&models.Task{}, id).Error
}
