package repository

import (
	"time"

	"github.com/Nehyan9895/reminder-system/internal/models"
)

func (r *GormRepo) CreateExecution(ruleID, taskID uint, triggeredAt time.Time) error {
	return r.DB.Create(&models.ReminderExecution{
		RuleID:      ruleID,
		TaskID:      taskID,
		TriggeredAt: triggeredAt,
	}).Error
}

func (r *GormRepo) HasExecutionSince(ruleID, taskID uint, since time.Time) (bool, error) {
	var cnt int64
	if err := r.DB.Model(&models.ReminderExecution{}).
		Where("rule_id = ? AND task_id = ? AND triggered_at >= ?", ruleID, taskID, since).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}
