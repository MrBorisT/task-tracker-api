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
				ID:   "11111111-1111-1111-1111-111111111111",
				Name: "Wake up",
			},
			{
				ID:   "22222222-2222-2222-2222-222222222222",
				Name: "Grab a brush",
			},
			{
				ID:   "33333333-3333-3333-3333-333333333333",
				Name: "Put a little make up",
			},
		},
	}

	r.Get("/health", app.HealthHandler)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", app.GetTasksHandler)
		r.Get("/{taskID}", app.GetTaskHandler)
		r.Delete("/{taskID}", app.DeleteTaskHandler)
		r.Post("/", app.CreateTaskHandler)
	})

	port := ":8080"
	log.Println("started server on port", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
