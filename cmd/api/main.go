package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/MrBorisT/task-tracker-api/internal/handlers"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func main() {
	pool, err := newPool()
	if err != nil {
		log.Fatalln("Unable to create database pool:", err)
	}

	defer pool.Close()
	r := chi.NewRouter()

	taskStore := storage.NewTaskStore(pool)

	r.Get("/health", handlers.HealthHandler(taskStore))
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handlers.GetTasksHandler(taskStore))
		r.Get("/{taskID}", handlers.GetTaskHandler(taskStore))
		r.Delete("/{taskID}", handlers.DeleteTaskHandler(taskStore))
		r.Post("/", handlers.CreateTaskHandler(taskStore))
		r.Put("/{taskID}", handlers.UpdateTaskHandler(taskStore))
	})

	port := ":8080"
	log.Println("started server on port", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func compileDSN() string {
	return "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " dbname=" + os.Getenv("DB_NAME") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=" + os.Getenv("DB_SSLMODE")
}

func newPool() (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	pool, err := pgxpool.New(context.Background(), compileDSN())
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
