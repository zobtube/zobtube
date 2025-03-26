package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusError      TaskStatus = "error"
)

type Task struct {
	ID         string `gorm:"type:uuid;primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	DoneAt     gorm.DeletedAt `gorm:"index"`
	Name       string
	Step       string
	Status     TaskStatus        `gorm:"default:todo"`
	Parameters map[string]string `gorm:"serializer:json"`
}

// UUID pre-hook
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "00000000-0000-0000-0000-000000000000" {
		t.ID = uuid.NewString()
		return nil
	}

	if t.ID == "" {
		t.ID = uuid.NewString()
		return nil
	}

	return nil
}
