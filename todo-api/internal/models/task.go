package models

import "github.com/google/uuid"

type Status string

const (
	StatusNotStarted Status = "notstarted"
	StatusOngoing    Status = "ongoing"
	StatusDone       Status = "done"
)

type Task struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
}

type TaskFilter struct {
	Status string
}

func IsValidStatus(status string) bool {
	switch Status(status) {
	case StatusNotStarted, StatusOngoing, StatusDone:
		return true
	default:
		return false
	}
}
