package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type Health struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	currentHealth := Health{
		Status: "ok",
	}
	//TODO check Health

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(currentHealth); err != nil {
		log.Println("encoding server health:", err)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	tasks := make([]Task, 0)
	if err := encoder.Encode(tasks); err != nil {
		log.Println("encoding tasks:", err)
	}
}

func PostTasks(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func main() {
	r := chi.NewRouter()

	r.Get("/health", HealthHandler)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", GetTasks)
		r.Post("/", PostTasks)
	})

	port := ":8080"
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalln(err)
	}
}
