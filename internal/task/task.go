package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusCreated   TaskStatus = "created"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID          uuid.UUID     `json:"id"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt time.Time  `json:"completed_at"`
	Result      string     `json:"result"`
}

func (t *Task) Duration() float64 {

	var endTime time.Time
	if t.CompletedAt.IsZero() { // rework logic //////////////////////////////////////////////////////////////////////////////////////////
		endTime = time.Now()
		t.CompletedAt = endTime;
	} else {			
		endTime = t.CompletedAt
	}

	return endTime.Sub(t.CreatedAt).Seconds()
}

func (t *Task) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":            t.ID,
		"status":        t.Status,
		"created_at":    t.CreatedAt.Format(time.RFC3339),
		"completed_at":  t.CompletedAt.Format(time.RFC3339),
		"duration_secs": t.Duration(),
		"result":        t.Result,
	}
}
