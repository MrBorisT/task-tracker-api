package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MrBorisT/task-tracker-api/handlers"
	"github.com/MrBorisT/task-tracker-api/models"
)

func main() {
	r := chi.NewRouter()
	app := handlers.App{
		Tasks: []models.Task{
			{
				ID:   0,
				Name: "Wake up",
				Done: false,
			},
			{
				ID:   1,
				Name: "Grab a brush",
				Done: false,
			},
			{
				ID:   2,
				Name: "Put a little make up",
				Done: false,
			},
		},
	}

	r.Get("/health", app.HealthHandler)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", app.GetTasksHandler)
		r.Post("/", app.CreateTaskHandler)
	})

	port := ":8080"
	log.Println("started server on port", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
