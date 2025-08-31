package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mayurbotre/task-management-system/internal/config"
	"github.com/mayurbotre/task-management-system/internal/repository"
	"github.com/mayurbotre/task-management-system/internal/service"
	th "github.com/mayurbotre/task-management-system/internal/transport/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load()

	dsn := cfg.MySQLDSN()
	if dsn == "" {
		log.Fatal("DATABASE_DSN or DB_* envs must be set for MySQL")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect MySQL: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db.DB(): %v", err)
	}
	tunePool(sqlDB)

	repo := repository.NewGormTaskRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	svc := service.NewTaskService(repo)
	router := th.SetupRouter(svc)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("task-service listening on %s (mysql)", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	} else {
		log.Println("server stopped")
	}

	log.Println("closing database pool...")
	if err := sqlDB.Close(); err != nil {
		log.Printf("db close error: %v", err)
	} else {
		log.Println("database pool closed")
	}
}

func tunePool(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
}
