package models

import "time"

// Task: seed at least 5 of these
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueAt       time.Time `json:"due_at"`
	Status      string    `json:"status"` // "pending", "done"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReminderRule: generic parameters encoded as JSON string (simple)
type ReminderRule struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Active    bool       `json:"active"`
	RuleType  string     `json:"rule_type"`               // "before_due", "interval","at_due"
	Params    string     `gorm:"type:TEXT" json:"params"` // JSON string
	LastRunAt *time.Time `json:"last_run_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// AuditLog stores actions and scheduler-triggered events
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventType string    `json:"event_type"` // rule.create, rule.update, reminder.trigger
	Details   string    `gorm:"type:TEXT" json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// ReminderExecution prevents duplicate triggers (one row per triggered rule+task)
type ReminderExecution struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RuleID      uint      `json:"rule_id"`
	TaskID      uint      `json:"task_id"`
	TriggeredAt time.Time `json:"triggered_at"`
	CreatedAt   time.Time `json:"created_at"`
}
