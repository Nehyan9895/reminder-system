package repository

import (
	"errors"
	"time"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"gorm.io/gorm"
)

func (r *GormRepo) CreateRule(rr *models.ReminderRule) error {
	return r.DB.Create(rr).Error
}

func (r *GormRepo) UpdateRule(rr *models.ReminderRule) error {
	return r.DB.Save(rr).Error
}

func (r *GormRepo) DeleteRule(id uint) error {
	return r.DB.Delete(&models.ReminderRule{}, id).Error
}

func (r *GormRepo) GetRuleByID(id uint) (*models.ReminderRule, error) {
	var rr models.ReminderRule
	if err := r.DB.First(&rr, id).Error; err != nil {
		return nil, err
	}
	return &rr, nil
}

func (r *GormRepo) ListRules() ([]models.ReminderRule, error) {
	var list []models.ReminderRule
	if err := r.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormRepo) ActiveRules() ([]models.ReminderRule, error) {
	var list []models.ReminderRule
	if err := r.DB.Where("active = ?", true).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormRepo) SetRuleActive(id uint, active bool) error {
	return r.DB.Model(&models.ReminderRule{}).
		Where("id = ?", id).
		Update("active", active).Error
}
func (r *GormRepo) LastExecutionTime(ruleID, taskID uint) (*time.Time, error) {
	var exec models.ReminderExecution
	err := r.DB.Where("rule_id = ? AND task_id = ?", ruleID, taskID).
		Order("triggered_at DESC").
		Limit(1).
		First(&exec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &exec.TriggeredAt, nil
}
