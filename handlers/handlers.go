package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

func (t *App) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskRequest := models.CreateTaskRequest{}

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	if err := decoder.Decode(&taskRequest); err != nil {
		log.Println("decoding task:", err)

		w.WriteHeader(http.StatusBadRequest)
		if err := encoder.Encode(ErrorResponse{Error: "invalid request body"}); err != nil {
			log.Println("encoding error response:", err)
		}
		return
	}

	trimmedName := strings.TrimSpace(taskRequest.Name)
	if trimmedName == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := encoder.Encode(ErrorResponse{Error: "task name cannot be empty"}); err != nil {
			log.Println("encoding error response:", err)
		}

		return
	}

	newTask := models.Task{
		Name: trimmedName,
		ID:   t.generateID(),
		Done: false,
	}

	t.Tasks = append(t.Tasks, newTask)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := encoder.Encode(newTask); err != nil {
		log.Println("post new task:", err)
	}
}

func (t *App) generateID() uint32 {
	return uuid.New().ID()
}
