package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/MrBorisT/task-tracker-api/internal/auth"
	"github.com/MrBorisT/task-tracker-api/internal/config"
	"github.com/MrBorisT/task-tracker-api/internal/handlers"
	"github.com/MrBorisT/task-tracker-api/internal/middleware"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading configuration:", err)
	}
	pool, err := newPool(config)
	if err != nil {
		log.Fatalln("Unable to create database pool:", err)
	}

	defer pool.Close()
	r := chi.NewRouter()

	taskStore := storage.NewTaskStore(pool)
	userStore := storage.NewUserStore(pool, config)
	authManager := auth.NewJWTManager(config)

	r.Get("/health", handlers.HealthHandler())
	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authManager))
		r.Get("/", handlers.GetTasksHandler(taskStore))
		r.Get("/{taskID}", handlers.GetTaskHandler(taskStore))
		r.Delete("/{taskID}", handlers.DeleteTaskHandler(taskStore))
		r.Post("/", handlers.CreateTaskHandler(taskStore))
		r.Put("/{taskID}", handlers.UpdateTaskHandler(taskStore))
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handlers.RegisterUserHandler(userStore))
		r.Post("/login", handlers.LoginUserHandler(userStore, authManager))
	})

	log.Println("started server on port", config.Port)
	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func compileDSN(config *config.Config) string {
	return "host=" + config.DBHost +
		" port=" + config.DBPort +
		" dbname=" + config.DBName +
		" user=" + config.DBUser +
		" password=" + config.DBPassword +
		" sslmode=" + config.DBSSLMode
}

func newPool(config *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dsn := compileDSN(config)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
