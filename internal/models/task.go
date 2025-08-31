package models

import "time"

type Status string

const (
	StatusPending    Status = "Pending"
	StatusInProgress Status = "InProgress"
	StatusCompleted  Status = "Completed"
)

func IsValidStatus(s Status) bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null;size:200"`
	Description string     `json:"description" gorm:"size:2000"`
	Status      Status     `json:"status" gorm:"type:text;default:Pending"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
