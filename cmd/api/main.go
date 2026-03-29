package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type Health struct {
	Status string `json:"status"`
}
type ErrorJson struct {
	Error string `json:"error"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	currentHealth := Health{
		Status: "ok",
	}
	//TODO check Health

	encoder := json.NewEncoder(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errAPI := ErrorJson{"method not allowed"}
		if err := encoder.Encode(errAPI); err != nil {
			log.Println("encoding error:", err)
		}
		return
	}
	if err := encoder.Encode(currentHealth); err != nil {
		log.Println("encoding server health:", err)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//TODO implement POST

	encoder := json.NewEncoder(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errAPI := ErrorJson{"method not allowed"}
		if err := encoder.Encode(errAPI); err != nil {
			log.Println("encoding error:", err)
		}
		return
	}
	tasks := make([]Task, 0)

	if err := encoder.Encode(tasks); err != nil {
		log.Println("encoding tasks:", err)
	}
}

func main() {
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/tasks", TasksHandler)

	port := ":8080"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
