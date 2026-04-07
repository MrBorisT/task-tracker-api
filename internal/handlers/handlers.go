package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func HealthHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func GetTasksHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		tasks, err := taskStore.ListTasks(r.Context())
		if err != nil {
			log.Println("listing tasks:", err)
			_ = writeJSONError(w, http.StatusInternalServerError, "error listing tasks")
			return
		}

		if err := encoder.Encode(tasks); err != nil {
			log.Println("encoding tasks:", err)
		}
	}
}

func GetTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))
		if task, err := taskStore.GetTask(r.Context(), taskID); err != nil {
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		} else if task != nil {
			if err := encoder.Encode(task); err != nil {
				log.Println("encoding task:", err)
			}
			return
		}
		_ = writeJSONError(w, http.StatusNotFound, "task not found")
	}
}

func DeleteTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

		// if err := validateTaskID(taskID); err != nil {
		// 	_ = writeJSONError(w, http.StatusBadRequest, "invalid task id")
		// 	return
		// }

		// taskStore.MU.Lock()
		// defer taskStore.MU.Unlock()
		// for i := range taskStore.Tasks {
		// 	if taskStore.Tasks[i].ID == taskID {
		// 		taskStore.Tasks = removeTaskByIndex(taskStore.Tasks, i)
		// 		w.WriteHeader(http.StatusNoContent)
		// 		return
		// 	}
		// }

		// _ = writeJSONError(w, http.StatusNotFound, "task not found")
	}
}

// func removeTaskByIndex(tasks []models.Task, index int) []models.Task {
// 	return append(tasks[:index], tasks[index+1:]...)
// }

func CreateTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// taskRequest := models.CreateTaskRequest{}

		// decoder := json.NewDecoder(r.Body)
		// encoder := json.NewEncoder(w)

		// w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// if err := decoder.Decode(&taskRequest); err != nil {
		// 	log.Println("decoding task:", err)
		// 	_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
		// 	return
		// }

		// taskStore.MU.Lock()
		// defer taskStore.MU.Unlock()
		// trimmedName := strings.TrimSpace(taskRequest.Name)
		// if trimmedName == "" {
		// 	_ = writeJSONError(w, http.StatusBadRequest, "task name cannot be empty")
		// 	return
		// }

		// newTask := models.Task{
		// 	Name:   trimmedName,
		// 	ID:     generateID(),
		// 	Status: models.StatusNew,
		// }

		// taskStore.Tasks = append(taskStore.Tasks, newTask)
		// w.WriteHeader(http.StatusCreated)

		// if err := encoder.Encode(newTask); err != nil {
		// 	log.Println("post new task:", err)
		// }
	}
}

func UpdateTaskHandler(taskStore *storage.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// taskID := strings.TrimSpace(chi.URLParam(r, "taskID"))

		// // if err := validateTaskID(taskID); err != nil {
		// // 	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		// // 	return
		// // }

		// taskRequest := models.UpdateTaskRequest{}

		// decoder := json.NewDecoder(r.Body)
		// encoder := json.NewEncoder(w)

		// w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// if err := decoder.Decode(&taskRequest); err != nil {
		// 	_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
		// 	return
		// }

		// taskStore.MU.Lock()

		// for i := range taskStore.Tasks {
		// 	if taskStore.Tasks[i].ID == taskID {
		// 		if err := updateTask(&taskStore.Tasks[i], taskRequest); err != nil {
		// 			taskStore.MU.Unlock()
		// 			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
		// 			return
		// 		}

		// 		taskStore.MU.Unlock()
		// 		w.WriteHeader(http.StatusOK)
		// 		if err := encoder.Encode(taskStore.Tasks[i]); err != nil {
		// 			log.Println("encoding updated task:", err)
		// 		}
		// 		return
		// 	}
		// }

		// taskStore.MU.Unlock()
		// _ = writeJSONError(w, http.StatusNotFound, "task not found")
	}
}

// get rid of these functions below!!!
// func updateTask(task *models.Task, taskRequest models.UpdateTaskRequest) error {
// 	if taskRequest.Name == nil && taskRequest.Status == nil {
// 		return fmt.Errorf("at least one field (name or status) must be provided for update")
// 	}
// 	if taskRequest.Name != nil {
// 		trimmedName := strings.TrimSpace(*taskRequest.Name)
// 		if trimmedName == "" {
// 			return fmt.Errorf("task name cannot be empty")
// 		}
// 		task.Name = trimmedName
// 	}
// 	if taskRequest.Status != nil {
// 		if !taskRequest.Status.IsValid() {
// 			return fmt.Errorf("invalid task status")
// 		}
// 		task.Status = *taskRequest.Status
// 	}
// 	return nil
// }

// func generateID() string {
// 	return uuid.New().String()
// }
