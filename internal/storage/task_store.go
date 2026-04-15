package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskStore struct {
	Pool *pgxpool.Pool
}

func NewTaskStore(pool *pgxpool.Pool) *TaskStore {
	return &TaskStore{Pool: pool}
}

func (s *TaskStore) ListTasks(ctx context.Context, userID string, gtq models.GetTasksQuery) ([]models.Task, error) {
	query := "SELECT id, name, status FROM tasks"
	args := []any{}
	argID := 1

	query += fmt.Sprintf(" WHERE user_id = $%d", argID)
	args = append(args, userID)
	argID++

	if gtq.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, gtq.Status)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d", argID)
	args = append(args, gtq.Limit)

	rows, err := s.Pool.Query(ctx, query, args...)
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

func (s *TaskStore) GetTask(ctx context.Context, userID string, id string) (*models.Task, error) {
	if !s.validateTaskID(id) {
		return nil, ErrInvalidTaskID
	}
	resultTask := models.Task{}

	query := "SELECT id, name, status FROM tasks WHERE user_id = $1 AND id = $2"
	if err := s.Pool.QueryRow(ctx, query, userID, id).Scan(&resultTask.ID, &resultTask.Name, &resultTask.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTaskNotFound
		} else {
			return nil, fmt.Errorf("error getting task: %w", err)
		}
	} else {
		return &resultTask, nil
	}
}

func (s *TaskStore) CreateTask(ctx context.Context, userID string, task models.CreateTaskRequest) (*models.Task, error) {
	trimmedName := strings.TrimSpace(task.Name)
	if trimmedName == "" {
		return nil, ErrEmptyTaskName
	}

	newTask := models.Task{
		ID:     s.generateID(),
		Name:   trimmedName,
		Status: models.StatusNew,
		UserID: userID,
	}

	query := "INSERT INTO tasks (id, name, status, user_id) VALUES ($1, $2, $3, $4) RETURNING id, name, status, user_id"
	row := s.Pool.QueryRow(ctx, query, newTask.ID, newTask.Name, newTask.Status, userID)

	if err := row.Scan(&newTask.ID, &newTask.Name, &newTask.Status, &newTask.UserID); err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	} else {
		return &newTask, nil
	}
}

func (s *TaskStore) DeleteTask(ctx context.Context, userID string, id string) error {
	if !s.validateTaskID(id) {
		return ErrInvalidTaskID
	}

	query := "DELETE FROM tasks WHERE user_id = $1 AND id = $2"
	commandTag, err := s.Pool.Exec(ctx, query, userID, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, userID string, task models.UpdateTaskRequest) (*models.Task, error) {
	if !s.validateTaskID(task.ID) {
		return nil, ErrInvalidTaskID
	}
	if task.Name == nil && task.Status == nil {
		return nil, ErrMissingUpdateFields
	}
	if task.Name != nil {
		if strings.TrimSpace(*task.Name) == "" {
			return nil, ErrEmptyTaskName
		}
	}
	if task.Status != nil {
		if !task.Status.IsValid() {
			return nil, ErrInvalidTaskStatus
		}
	}

	query := "UPDATE tasks SET name = COALESCE($1, name), status = COALESCE($2, status) WHERE user_id = $3 AND id = $4 RETURNING id, name, status"

	updatedTask := models.Task{ID: task.ID}
	row := s.Pool.QueryRow(ctx, query, task.Name, task.Status, userID, task.ID)
	if err := row.Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTaskNotFound
		} else {
			return nil, fmt.Errorf("error updating task: %w", err)
		}
	}
	return &updatedTask, nil
}

func (s *TaskStore) validateTaskID(taskID string) bool {
	_, err := uuid.Parse(taskID)
	if err != nil {
		return false
	}
	return true
}

func (s *TaskStore) generateID() string {
	return uuid.New().String()
}
