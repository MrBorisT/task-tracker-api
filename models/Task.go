package models

type Task struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}
