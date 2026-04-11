package storage

import "errors"

var (
	//task
	ErrInvalidTaskID       = errors.New("invalid task ID")
	ErrTaskNotFound        = errors.New("task not found")
	ErrEmptyTaskName       = errors.New("task name cannot be empty")
	ErrInvalidTaskStatus   = errors.New("invalid task status")
	ErrMissingUpdateFields = errors.New("at least one field (name or status) must be provided for update")

	//user
	ErrEmptyUserEmail    = errors.New("user email cannot be empty")
	ErrEmptyUserPassword = errors.New("user password cannot be empty")
	ErrShortUserPassword = errors.New("user password must be at least 6 characters long")
	ErrLongUserPassword  = errors.New("user password cannot be longer than 72 characters")
)
