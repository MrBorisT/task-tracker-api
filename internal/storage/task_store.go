package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskStore struct {
	MU   sync.RWMutex
	Pool *pgxpool.Pool
}

func (s *TaskStore) ListTasks(ctx context.Context) ([]models.Task, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id, name, status FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultTasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Status); err != nil {
			return nil, err
		}
		resultTasks = append(resultTasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultTasks, nil
}

func (s *TaskStore) GetTask(ctx context.Context, id string) (*models.Task, error) {
	if err := s.validateTaskID(id); err != nil {
		return nil, err
	}
	resultTask := models.Task{}

	if err := s.Pool.QueryRow(ctx, "SELECT id, name, status FROM tasks WHERE id = $1", id).Scan(&resultTask.ID, &resultTask.Name, &resultTask.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		} else {
			return nil, fmt.Errorf("error getting task: %w", err)
		}
	} else {
		return &resultTask, nil
	}
}

func (s *TaskStore) validateTaskID(taskID string) error {
	_, err := uuid.Parse(taskID)
	if err != nil {
		return fmt.Errorf("invalid task id")
	}
	return nil
}
