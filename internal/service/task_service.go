package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mayurbotre/task-management-system/internal/models"
	"github.com/mayurbotre/task-management-system/internal/repository"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(ctx context.Context, in CreateTaskInput) (*models.Task, error)
	GetTask(ctx context.Context, id uint) (*models.Task, error)
	UpdateTask(ctx context.Context, id uint, in UpdateTaskInput) (*models.Task, error)
	DeleteTask(ctx context.Context, id uint) error
	ListTasks(ctx context.Context, f ListFilter) ([]models.Task, int64, error)
}

type taskService struct{ repo repository.TaskRepository }

func NewTaskService(repo repository.TaskRepository) TaskService { return &taskService{repo: repo} }

type CreateTaskInput struct {
	Title       string
	Description string
	Status      *models.Status
	DueDate     *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *models.Status
	DueDate     *time.Time
}

type ListFilter struct {
	Status   *models.Status
	Page     int
	PageSize int
}

func (s *taskService) CreateTask(ctx context.Context, in CreateTaskInput) (*models.Task, error) {
	if in.Title == "" {
		return nil, errors.New("title is required")
	}
	task := &models.Task{
		Title:       in.Title,
		Description: in.Description,
		Status:      models.StatusPending,
		DueDate:     in.DueDate,
	}
	if in.Status != nil {
		if !models.IsValidStatus(*in.Status) {
			return nil, fmt.Errorf("invalid status: %s", *in.Status)
		}
		task.Status = *in.Status
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, id uint) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) UpdateTask(ctx context.Context, id uint, in UpdateTaskInput) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Title != nil {
		if *in.Title == "" {
			return nil, errors.New("title cannot be empty")
		}
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = *in.Description
	}
	if in.Status != nil {
		if !models.IsValidStatus(*in.Status) {
			return nil, fmt.Errorf("invalid status: %s", *in.Status)
		}
		task.Status = *in.Status
	}
	if in.DueDate != nil {
		task.DueDate = in.DueDate
	}
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *taskService) ListTasks(ctx context.Context, f ListFilter) ([]models.Task, int64, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 10
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	return s.repo.List(ctx, repository.TaskFilter{
		Status: f.Status,
		Pagination: repository.Pagination{
			Page:     f.Page,
			PageSize: f.PageSize,
		},
	})
}
