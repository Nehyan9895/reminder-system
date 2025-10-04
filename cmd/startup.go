package main

import (
	"time"

	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
)

func seedIfEmpty(repo *repository.GormRepo) {
	var cnt int64
	repo.DB.Model(&models.Task{}).Count(&cnt)
	if cnt > 0 {
		return
	}

	now := time.Now()
	tasks := []models.Task{
		{Title: "Pay electricity bill", Description: "Electricity", DueAt: now.Add(2 * time.Minute), Status: "pending"},
		{Title: "Submit assignment", Description: "Bootcamp", DueAt: now.Add(10 * time.Minute), Status: "pending"},
		{Title: "Daily workout", Description: "Run", DueAt: now.Add(1 * time.Hour), Status: "pending"},
		{Title: "Call supplier", Description: "Discuss order", DueAt: now.Add(3 * time.Minute), Status: "pending"},
		{Title: "Read chapter 4", Description: "Study", DueAt: now.Add(20 * time.Minute), Status: "pending"},
	}
	for i := range tasks {
		_ = repo.CreateTask(&tasks[i])
	}

	// create sample rules
	bparams := `{"minutes_before":1}`
	iparams := `{"interval_min":2}`
	_ = repo.CreateRule(&models.ReminderRule{Name: "1min before", Active: true, RuleType: "before_due", Params: bparams})
	_ = repo.CreateRule(&models.ReminderRule{Name: "every 2 min", Active: true, RuleType: "interval", Params: iparams})
	_ = repo.WriteAudit("seed", "seeded sample tasks and rules")
}
