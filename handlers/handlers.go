package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MrBorisT/task-tracker-api/models"
)

type App struct {
	Tasks []models.Task
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

	if err := encoder.Encode(t.Tasks); err != nil {
		log.Println("encoding tasks:", err)
	}
}

func (t *App) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)

	taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

	_, err := uuid.Parse(taskID)
	if err != nil {
		log.Println("parsing task id", err)
		writeJSONError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	for _, task := range t.Tasks {
		if task.ID == taskID {
			if err := encoder.Encode(task); err != nil {
				log.Println("encoding task:", err)
				if err = writeJSONError(w, http.StatusBadRequest, "failed to encode task"); err != nil {
					log.Println("encoding error response:", err)
				}
			}
			return
		}
	}

	if err := writeJSONError(w, http.StatusNotFound, "task not found"); err != nil {
		log.Println("encoding error response:", err)
	}
}

func (t *App) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskRequest := models.CreateTaskRequest{}

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := decoder.Decode(&taskRequest); err != nil {
		log.Println("decoding task:", err)

		if err := writeJSONError(w, http.StatusBadRequest, "invalid request body"); err != nil {
			log.Println("encoding error response:", err)
		}
		return
	}

	trimmedName := strings.TrimSpace(taskRequest.Name)
	if trimmedName == "" {
		if err := writeJSONError(w, http.StatusBadRequest, "task name cannot be empty"); err != nil {
			log.Println("encoding error response:", err)
		}

		return
	}

	newTask := models.Task{
		Name: trimmedName,
		ID:   generateID(),
		Done: false,
	}

	t.Tasks = append(t.Tasks, newTask)
	w.WriteHeader(http.StatusCreated)

	if err := encoder.Encode(newTask); err != nil {
		log.Println("post new task:", err)
	}
}

func generateID() string {
	return uuid.New().String()
}
