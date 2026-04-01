package models

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}
