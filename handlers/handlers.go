package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
	if err := decoder.Decode(&taskRequest); err != nil {
		log.Println("decoding task:", err)

		//todo return json error
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if taskRequest.Name == "" {

		//todo return json error
		http.Error(w, "task name is required", http.StatusBadRequest)
		return
	}

	newTask := models.Task{
		Name: taskRequest.Name,
		ID:   t.generateID(),
		Done: false,
	}

	t.Tasks = append(t.Tasks, newTask)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(newTask); err != nil {
		log.Println("post new task:", err)
	}
}

func (t *App) generateID() int {
	return len(t.Tasks)
}
