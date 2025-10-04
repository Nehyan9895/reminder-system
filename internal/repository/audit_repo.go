package repository

import "github.com/Nehyan9895/reminder-system/internal/models"

func (r *GormRepo) WriteAudit(eventType, details string) error {
	return r.DB.Create(&models.AuditLog{
		EventType: eventType,
		Details:   details,
	}).Error
}
