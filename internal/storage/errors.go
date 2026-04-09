package storage

import "errors"

var (
	ErrInvalidTaskID       = errors.New("invalid task ID")
	ErrTaskNotFound        = errors.New("task not found")
	ErrEmptyTaskName       = errors.New("task name cannot be empty")
	ErrInvalidTaskStatus   = errors.New("invalid task status")
	ErrMissingUpdateFields = errors.New("at least one field (name or status) must be provided for update")
)
