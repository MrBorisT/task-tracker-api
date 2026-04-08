package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func HealthHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		currentHealth := models.Health{
			Status: "ok",
		}
		//TODO check Health

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(currentHealth); err != nil {
			log.Println("encoding server health:", err)
		}
	}
}

func GetTasksHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		tasks, err := taskStore.ListTasks(r.Context())
		if err != nil {
			log.Println("listing tasks:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "error listing tasks")
			return
		}

		if err := encoder.Encode(tasks); err != nil {
			log.Println("encoding tasks:", err)
		}
	}
}

func GetTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))
		if task, err := taskStore.GetTask(r.Context(), taskID); err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		} else if task != nil {
			if err := encoder.Encode(task); err != nil {
				log.Println("encoding task:", err)
			}
			return
		}
		_ = writeJSONError(w, http.StatusNotFound, "task not found")
	}
}

func CreateTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskRequest := models.CreateTaskRequest{}

		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&taskRequest); err != nil {
			log.Println("decoding task:", err)
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		newTask, err := taskStore.CreateTask(r.Context(), taskRequest)
		if err != nil {
			log.Println("creating task:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "error creating task")
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := encoder.Encode(newTask); err != nil {
			log.Println("post new task:", err)
		}
	}
}

func DeleteTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

		if err := taskStore.DeleteTask(r.Context(), taskID); err != nil {
			log.Println("deleting task:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "error deleting task")
			return
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
			log.Println("updating task:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "error updating task")
			return
		} else if newTask != nil {
			if err := encoder.Encode(newTask); err != nil {
				log.Println("encoding updated task:", err)
			}
			return
		}
		_ = writeJSONError(w, http.StatusNotFound, "task not found")
	}
}
