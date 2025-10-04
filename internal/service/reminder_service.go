package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
	log "github.com/sirupsen/logrus"
)

type ReminderService struct {
	repo *repository.GormRepo
}

func NewReminderService(r *repository.GormRepo) *ReminderService {
	return &ReminderService{repo: r}
}

type BeforeDueParams struct {
	MinutesBefore int `json:"minutes_before"`
}
type IntervalParams struct {
	IntervalMin int `json:"interval_min"`
}

// RunOnce: applies rules one pass (idempotent via executions)
func (s *ReminderService) RunOnce() {
	log.Info("[scheduler] run pass")
	rules, err := s.repo.ActiveRules()
	if err != nil {
		log.Errorf("fetch rules: %v", err)
		return
	}

	for _, rr := range rules {
		switch rr.RuleType {
		case "before_due":
			s.applyBeforeDue(&rr)
		case "interval":
			s.applyInterval(&rr)
		case "at_due":
			s.applyAtDue(&rr)
		default:
			log.Warnf("unknown rule type: %s", rr.RuleType)
		}
		now := time.Now()
		rr.LastRunAt = &now
		_ = s.repo.UpdateRule(&rr)
	}
}

func (s *ReminderService) applyBeforeDue(rr *models.ReminderRule) {
	var p BeforeDueParams
	if err := json.Unmarshal([]byte(rr.Params), &p); err != nil {
		log.Errorf("invalid params for rule %d: %v", rr.ID, err)
		return
	}

	now := time.Now()

	windowStart := now
	windowEnd := now.Add(time.Duration(p.MinutesBefore) * time.Minute)

	var tasks []models.Task
	if err := s.repo.DB.
		Where("status = ? AND due_at BETWEEN ? AND ?", "pending", windowStart, windowEnd).
		Find(&tasks).Error; err != nil {
		log.Errorf("fetch tasks: %v", err)
		return
	}

	for _, t := range tasks {
		already, err := s.repo.HasExecutionSince(rr.ID, t.ID, windowStart)
		if err != nil {
			log.Errorf("check exec: %v", err)
			continue
		}
		if already {
			continue
		}

		msg := fmt.Sprintf("Reminder(rule:%s) -> Task:%d %s due:%s", rr.Name, t.ID, t.Title, t.DueAt.Format(time.RFC3339))
		log.Info(msg)

		details := fmt.Sprintf(
			"Reminder triggered [Rule #%d: %s] -> [Task #%d: %s]",
			rr.ID, rr.Name, t.ID, t.Title,
		)
		_ = s.repo.WriteAudit("reminder.trigger", details)
		_ = s.repo.CreateExecution(rr.ID, t.ID, now)
	}
}

func (s *ReminderService) applyInterval(rr *models.ReminderRule) {
	var p IntervalParams
	if err := json.Unmarshal([]byte(rr.Params), &p); err != nil {
		log.Errorf("invalid params for rule %d: %v", rr.ID, err)
		return
	}

	now := time.Now()
	var tasks []models.Task
	// Get only tasks past due and still pending
	if err := s.repo.DB.Where("status = ? AND due_at <= ?", "pending", now).Find(&tasks).Error; err != nil {
		log.Errorf("fetch tasks: %v", err)
		return
	}

	for _, t := range tasks {
		// Check last execution
		lastExec, err := s.repo.LastExecutionTime(rr.ID, t.ID)
		if err != nil {
			log.Errorf("get last exec: %v", err)
			continue
		}

		// If never reminded OR enough time has passed
		if lastExec == nil || now.Sub(*lastExec) >= time.Duration(p.IntervalMin)*time.Minute {
			msg := fmt.Sprintf("IntervalReminder(rule:%s) -> Task:%d %s (past due: %s)",
				rr.Name, t.ID, t.Title, t.DueAt.Format(time.RFC3339))
			log.Info(msg)

			details := fmt.Sprintf(
				"Reminder triggered [Rule #%d: %s] -> [Task #%d: %s]",
				rr.ID, rr.Name, t.ID, t.Title,
			)

			_ = s.repo.WriteAudit("reminder.trigger", details)
			_ = s.repo.CreateExecution(rr.ID, t.ID, time.Now())
		}
	}
}

func (s *ReminderService) applyAtDue(rr *models.ReminderRule) {
	// Parse params if needed (could include window tolerance, e.g., 1 min)
	now := time.Now()
	var tasks []models.Task
	if err := s.repo.DB.Where("status = ? AND due_at BETWEEN ? AND ?", "pending", now.Add(-1*time.Minute), now.Add(1*time.Minute)).Find(&tasks).Error; err != nil {
		log.Errorf("fetch tasks: %v", err)
		return
	}

	for _, t := range tasks {
		lastExec, err := s.repo.LastExecutionTime(rr.ID, t.ID)
		if err != nil {
			log.Errorf("get last exec: %v", err)

			continue
		}
		if lastExec != nil {
			continue // already reminded
		}

		if now.Equal(t.DueAt) || now.After(t.DueAt) {
			msg := fmt.Sprintf("AtDueReminder(rule:%s) -> Task:%d %s due:%s", rr.Name, t.ID, t.Title, t.DueAt.Format(time.RFC3339))
			log.Info(msg)

			details := fmt.Sprintf(
				"Reminder triggered [Rule #%d: %s] -> [Task #%d: %s]",
				rr.ID, rr.Name, t.ID, t.Title,
			)

			_ = s.repo.WriteAudit("reminder.trigger", details)
			_ = s.repo.CreateExecution(rr.ID, t.ID, now)
		}

	}
}

// StartScheduler runs periodic loop in a goroutine and returns a cancel function via context
func (s *ReminderService) StartScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	// run once immediately
	s.RunOnce()
	for {
		select {
		case <-ctx.Done():
			log.Info("scheduler stopping")
			return
		case <-ticker.C:
			s.RunOnce()
		}
	}
}
