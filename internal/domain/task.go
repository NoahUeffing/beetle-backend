package domain

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in the system
// @Description A task entity with its properties
type Task struct {
	// The unique identifier of the task
	ID uint `json:"id" gorm:"primaryKey" example:"1"`
	// The title of the task
	Title string `json:"title" gorm:"not null" example:"Complete project documentation"`
	// A detailed description of the task
	Description string `json:"description" example:"Write comprehensive documentation for the API endpoints"`
	// The current status of the task (pending, in_progress, completed)
	Status string `json:"status" gorm:"default:'pending'" example:"pending"`
	// The timestamp when the task was created
	CreatedAt time.Time `json:"created_at" example:"2024-02-27T12:00:00Z"`
	// The timestamp when the task was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2024-02-27T12:00:00Z"`
	// The timestamp when the task was deleted (soft delete)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time" example:"2024-02-27T12:00:00Z"`
}

func (Task) TableName() string {
	return "tasks"
}
