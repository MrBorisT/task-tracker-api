package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MrBorisT/task-tracker-api/models"
)

type App struct {
	Tasks []models.Task
	mu    sync.RWMutex
}

func (t *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
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

func (t *App) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)

	t.mu.RLock()
	tasksCopy := make([]models.Task, len(t.Tasks))
	copy(tasksCopy, t.Tasks)
	t.mu.RUnlock()

	if err := encoder.Encode(tasksCopy); err != nil {
		log.Println("encoding tasks:", err)
	}
}

func (t *App) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)

	taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

	t.mu.RLock()
	defer t.mu.RUnlock()
	_, err := uuid.Parse(taskID)
	if err != nil {
		log.Println("parsing task id", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	for _, task := range t.Tasks {
		if task.ID == taskID {
			if err = encoder.Encode(task); err != nil {
				log.Println("encoding task:", err)
			}
			return
		}
	}

	_ = writeJSONError(w, http.StatusNotFound, "task not found")
}

func (t *App) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskRequest := models.CreateTaskRequest{}

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := decoder.Decode(&taskRequest); err != nil {
		log.Println("decoding task:", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	trimmedName := strings.TrimSpace(taskRequest.Name)
	if trimmedName == "" {
		_ = writeJSONError(w, http.StatusBadRequest, "task name cannot be empty")
		return
	}

	newTask := models.Task{
		Name: trimmedName,
		ID:   t.generateID(),
		Done: false,
	}

	t.Tasks = append(t.Tasks, newTask)
	w.WriteHeader(http.StatusCreated)

	if err := encoder.Encode(newTask); err != nil {
		log.Println("post new task:", err)
	}
}

func (t *App) generateID() string {
	return uuid.New().String()
}
