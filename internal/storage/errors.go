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
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

const (
	PGCodeUniqueViolation = "23505"
)
