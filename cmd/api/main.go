package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/MrBorisT/task-tracker-api/handlers"
	"github.com/MrBorisT/task-tracker-api/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " dbname=" + os.Getenv("DB_NAME") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=" + os.Getenv("DB_SSLMODE")
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Println("Unable to connect to database:", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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
		r.Put("/{taskID}", app.UpdateTaskHandler)
	})

	port := ":8080"
	log.Println("started server on port", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
