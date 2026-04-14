package models

type Task struct {
	ID     string     `json:"id,omitempty"`
	Name   string     `json:"name,omitempty"`
	Status TaskStatus `json:"status"`
	UserID string     `json:"user_id,omitempty"`
}

type CreateTaskRequest struct {
	Name string `json:"name,omitempty"`
}

type UpdateTaskRequest struct {
	ID     string      `json:"id,omitempty"`
	Name   *string     `json:"name,omitempty"`
	Status *TaskStatus `json:"status,omitempty"`
}

type GetTasksQuery struct {
	Status string
	Limit  int
}
