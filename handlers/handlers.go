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
	// newTask := Task{}

}
