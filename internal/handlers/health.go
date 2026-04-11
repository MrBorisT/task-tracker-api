package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MrBorisT/task-tracker-api/internal/models"
)

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		currentHealth := models.Health{
			Status: "ok",
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(currentHealth); err != nil {
			log.Println("encoding server health:", err)
		}
	}
}
