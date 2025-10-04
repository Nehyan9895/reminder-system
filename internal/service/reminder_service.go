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
	// find pending tasks with due_at <= now + minutesBefore and due_at > now - 1h (guard)
	now := time.Now()
	threshold := now.Add(time.Duration(p.MinutesBefore) * time.Minute)
	var tasks []models.Task
	// raw GORM query to fetch tasks
	if err := s.repo.DB.Where("status = ? AND due_at <= ? AND due_at >= ?", "pending", threshold, now.Add(-24*time.Hour)).Find(&tasks).Error; err != nil {
		log.Errorf("fetch tasks: %v", err)
		return
	}

	for _, t := range tasks {
		// dedupe: don't trigger multiple times within the same window (e.g., since due - 1day)
		windowStart := t.DueAt.Add(-time.Duration(p.MinutesBefore) * time.Minute).Add(-1 * time.Minute)
		already, err := s.repo.HasExecutionSince(rr.ID, t.ID, windowStart)
		if err != nil {
			log.Errorf("check exec: %v", err)
			continue
		}
		if already {
			continue
		}

		// simulate send
		msg := fmt.Sprintf("Reminder(rule:%s) -> Task:%d %s due:%s", rr.Name, t.ID, t.Title, t.DueAt.Format(time.RFC3339))
		log.Info(msg)

		// write audit
		details := fmt.Sprintf(`{"rule_id":%d,"rule_name":"%s","task_id":%d,"task_title":"%s"}`, rr.ID, rr.Name, t.ID, t.Title)
		_ = s.repo.WriteAudit("reminder.trigger", details)

		// record execution
		_ = s.repo.CreateExecution(rr.ID, t.ID, time.Now())
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

			details := fmt.Sprintf(`{"rule_id":%d,"rule_name":"%s","task_id":%d}`, rr.ID, rr.Name, t.ID)
			_ = s.repo.WriteAudit("reminder.trigger", details)
			_ = s.repo.CreateExecution(rr.ID, t.ID, time.Now())
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
