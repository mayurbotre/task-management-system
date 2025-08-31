package repository

import (
	"context"

	"github.com/mayurbotre/task-management-system/internal/models"
)

func (r *gormTaskRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&models.Task{})
}

func (r *gormTaskRepository) Create(ctx context.Context, t *models.Task) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *gormTaskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *gormTaskRepository) Update(ctx context.Context, t *models.Task) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *gormTaskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

func (r *gormTaskRepository) List(ctx context.Context, f TaskFilter) ([]models.Task, int64, error) {
	var tasks []models.Task
	q := r.db.WithContext(ctx).Model(&models.Task{})
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (f.Page - 1) * f.PageSize
	if offset < 0 {
		offset = 0
	}
	if err := q.Order("created_at DESC").Limit(f.PageSize).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}
