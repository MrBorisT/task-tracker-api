package models

type Task struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Done bool   `json:"done"`
}

type CreateTaskRequest struct {
	Name string `json:"name,omitempty"`
}

type UpdateTaskRequest struct {
	Name string `json:"name,omitempty"`
	Done bool   `json:"done,omitempty"`
}
