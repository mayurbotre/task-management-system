package http

import (
	"time"

	"github.com/mayurbotre/task-management-system/internal/models"
)

type createTaskRequest struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	Status      *models.Status `json:"status,omitempty"`
	DueDate     *time.Time     `json:"dueDate,omitempty"`
}

type updateTaskRequest struct {
	Title       *string        `json:"title,omitempty"`
	Description *string        `json:"description,omitempty"`
	Status      *models.Status `json:"status,omitempty"`
	DueDate     *time.Time     `json:"dueDate,omitempty"`
}

type listTaskResponse struct {
	Items []models.Task  `json:"items"`
	Meta  map[string]any `json:"meta"`
}
