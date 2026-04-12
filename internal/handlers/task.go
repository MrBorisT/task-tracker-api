package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func GetTasksHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		q := r.URL.Query()
		query, err := newGetTasksQuery(q.Get("status"), q.Get("limit"))
		if err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, "invalid query parameters")
			return
		}

		tasks, err := taskStore.ListTasks(r.Context(), *query)
		if err != nil {
			log.Println("listing tasks:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		if err := encoder.Encode(tasks); err != nil {
			log.Println("encoding tasks:", err)
		}
	}
}

func newGetTasksQuery(statusStr, limitStr string) (*models.GetTasksQuery, error) {
	status := models.TaskStatus(statusStr)
	if status != "" {
		if !status.IsValid() {
			return nil, errors.New("invalid task status: " + statusStr)
		}
	}
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, errors.New("invalid limit: " + limitStr)
	}
	if limit < 0 {
		return nil, errors.New("limit cannot be negative: " + limitStr)
	}
	query := &models.GetTasksQuery{
		Status: string(status),
		Limit:  limit,
	}
	return query, nil
}

func GetTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))
		if task, err := taskStore.GetTask(r.Context(), taskID); err != nil {
			switch {
			case errors.Is(err, storage.ErrInvalidTaskID):
				_ = writeJSONError(w, http.StatusBadRequest, "invalid task ID")
				return
			case errors.Is(err, storage.ErrTaskNotFound):
				_ = writeJSONError(w, http.StatusNotFound, "task not found")
				return
			}
			log.Println("getting task:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		} else if task != nil {
			if err := encoder.Encode(task); err != nil {
				log.Println("encoding task:", err)
			}
			return
		}
	}
}

func CreateTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskRequest := models.CreateTaskRequest{}

		decoder := json.NewDecoder(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&taskRequest); err != nil {
			log.Println("decoding task:", err)
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newTask, err := taskStore.CreateTask(r.Context(), taskRequest)
		if err != nil {
			if err == storage.ErrEmptyTaskName {
				_ = writeJSONError(w, http.StatusBadRequest, "task name cannot be empty")
				return
			}
			log.Println("creating task:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		w.WriteHeader(http.StatusCreated)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(newTask); err != nil {
			log.Println("post new task:", err)
		}
	}
}

func DeleteTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

		if err := taskStore.DeleteTask(r.Context(), taskID); err != nil {
			switch {
			case errors.Is(err, storage.ErrInvalidTaskID):
				_ = writeJSONError(w, http.StatusBadRequest, "invalid task ID")
				return
			case errors.Is(err, storage.ErrTaskNotFound):
				_ = writeJSONError(w, http.StatusNotFound, "task not found")
				return
			default:
				log.Println("deleting task:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
				return
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func UpdateTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))
		taskRequest := models.UpdateTaskRequest{ID: taskID}

		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&taskRequest); err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if newTask, err := taskStore.UpdateTask(r.Context(), taskRequest); err != nil {
			switch {
			case errors.Is(err, storage.ErrInvalidTaskID):
				_ = writeJSONError(w, http.StatusBadRequest, "invalid task ID")
				return
			case errors.Is(err, storage.ErrEmptyTaskName):
				_ = writeJSONError(w, http.StatusBadRequest, "task name cannot be empty")
				return
			case errors.Is(err, storage.ErrInvalidTaskStatus):
				_ = writeJSONError(w, http.StatusBadRequest, "invalid task status")
				return
			case errors.Is(err, storage.ErrMissingUpdateFields):
				_ = writeJSONError(w, http.StatusBadRequest, "at least one field (name or status) must be provided for update")
				return
			default:
				log.Println("updating task:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
				return
			}
		} else if newTask != nil {
			if err := encoder.Encode(newTask); err != nil {
				log.Println("encoding updated task:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			}
			return
		}
	}
}
