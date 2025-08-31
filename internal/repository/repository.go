package repository

import (
	"context"

	"github.com/mayurbotre/task-management-system/internal/models"
	"gorm.io/gorm"
)

type Pagination struct {
	Page     int
	PageSize int
}

type TaskFilter struct {
	Status *models.Status
	Pagination
}

type TaskRepository interface {
	AutoMigrate() error
	Create(ctx context.Context, t *models.Task) error
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, f TaskFilter) ([]models.Task, int64, error)
}

type gormTaskRepository struct{ db *gorm.DB }

func NewGormTaskRepository(db *gorm.DB) TaskRepository { return &gormTaskRepository{db: db} }
