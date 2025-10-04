package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Nehyan9895/reminder-system/config"
	"github.com/Nehyan9895/reminder-system/internal/handler"
	"github.com/Nehyan9895/reminder-system/internal/models"
	"github.com/Nehyan9895/reminder-system/internal/repository"
	"github.com/Nehyan9895/reminder-system/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// DB
	config.LoadEnv()
	dsn := config.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	// Automigrate
	if err := db.AutoMigrate(&models.Task{}, &models.ReminderRule{}, &models.AuditLog{}, &models.ReminderExecution{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Repos
	repo := repository.NewGormRepo(db)

	// Services
	reminderSvc := service.NewReminderService(repo)
	taskSvc := service.NewTaskService(repo)

	// Handlers
	reminderHandler := handler.NewReminderHandler(repo)
	auditHandler := handler.NewAuditHandler(repo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	// Router
	r := chi.NewRouter()

	// Register routes
	reminderHandler.Register(r)
	auditHandler.Register(r)
	taskHandler.Register(r)

	// Seed sample tasks & rules
	seedIfEmpty(repo)

	// Scheduler
	ctx, cancel := context.WithCancel(context.Background())
	go reminderSvc.StartScheduler(ctx, 60*time.Second)

	// Serve UI static files
	r.Handle("/*", http.FileServer(http.Dir("./ui")))

	server := &http.Server{Addr: ":" + config.HTTPPort(), Handler: r}

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("shutting down...")
		cancel()
		_ = server.Shutdown(context.Background())
	}()

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
