package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func RegisterUserHandler(userStore *storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.RegisterUserRequest{}

		decoder := json.NewDecoder(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := userStore.RegisterUser(r.Context(), userRequest); err != nil {
			switch err {
			case storage.ErrEmptyUserEmail:
				_ = writeJSONError(w, http.StatusBadRequest, "user email cannot be empty")
			case storage.ErrEmptyUserPassword:
				_ = writeJSONError(w, http.StatusBadRequest, "user password cannot be empty")
			case storage.ErrShortUserPassword:
				_ = writeJSONError(w, http.StatusBadRequest, "user password must be at least 6 characters long")
			case storage.ErrLongUserPassword:
				_ = writeJSONError(w, http.StatusBadRequest, "user password cannot be longer than 72 characters")
			default:
				log.Println("registering user:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "error registering user")
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
