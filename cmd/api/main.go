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
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " dbname=" + os.Getenv("DB_NAME") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=" + os.Getenv("DB_SSLMODE")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Println("Unable to connect to database:", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Unable to ping database:", err)
		os.Exit(1)
	}

	//testing select statement
	// rows, err := pool.Query(context.Background(), "SELECT id, name, status FROM tasks")
	// if err != nil {
	// 	log.Println("Unable to execute query:", err)
	// 	os.Exit(1)
	// }
	// for rows.Next() {
	// 	newTask := models.Task{}
	// 	if err := rows.Scan(&newTask.ID, &newTask.Name, &newTask.Status); err != nil {
	// 		log.Println("Unable to scan row:", err)
	// 		os.Exit(1)
	// 	}
	// 	log.Println("Task: ", newTask)
	// }

	r := chi.NewRouter()

	taskStore := storage.TaskStore{
		Pool: pool,
	}

	r.Get("/health", handlers.HealthHandler(&taskStore))
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handlers.GetTasksHandler(&taskStore))
		r.Get("/{taskID}", handlers.GetTaskHandler(&taskStore))
		r.Delete("/{taskID}", handlers.DeleteTaskHandler(&taskStore))
		r.Post("/", handlers.CreateTaskHandler(&taskStore))
		r.Put("/{taskID}", handlers.UpdateTaskHandler(&taskStore))
	})

	port := ":8080"
	log.Println("started server on port", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
