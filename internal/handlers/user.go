package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/MrBorisT/task-tracker-api/internal/models"
	"github.com/MrBorisT/task-tracker-api/internal/storage"
)

func RegisterUserHandler(userStore *storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.UserRequest{}

		decoder := json.NewDecoder(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := verifyUserRequest(&userRequest); err != nil {
			log.Println("verifying user request:", err)
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := userStore.RegisterUser(r.Context(), userRequest); err != nil {
			switch err {
			case storage.ErrUserAlreadyExists:
				_ = writeJSONError(w, http.StatusConflict, "user with this email already exists")
			default:
				log.Println("registering user:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "error registering user")
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func LoginUserHandler(userStore *storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.UserRequest{}

		decoder := json.NewDecoder(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := verifyUserRequest(&userRequest); err != nil {
			log.Println("verifying user request:", err)
			_ = writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		userID, err := userStore.GetUserID(r.Context(), userRequest)
		if err != nil {
			switch err {
			case storage.ErrUserNotFound:
				_ = writeJSONError(w, http.StatusNotFound, "user not found")
			case storage.ErrInvalidCredentials:
				_ = writeJSONError(w, http.StatusUnauthorized, "invalid email or password")
			default:
				log.Println("logging in user:", err)
				_ = writeJSONError(w, http.StatusInternalServerError, "error logging in user")
			}
			return
		}
		//todo generate jwt and return it!
		fmt.Println(userID)
	}
}

func verifyUserRequest(userRequest *models.UserRequest) error {
	trimmedEmail := strings.TrimSpace(userRequest.Email)
	if trimmedEmail == "" {
		return ErrEmptyUserEmail
	}
	if _, err := mail.ParseAddress(trimmedEmail); err != nil {
		return ErrIncorrectUserEmail
	}

	trimmedPassword := strings.TrimSpace(userRequest.Password)
	if trimmedPassword == "" {
		return ErrEmptyUserPassword
	}
	if len(trimmedPassword) < 6 {
		return ErrShortUserPassword
	}
	if len(trimmedPassword) > 72 {
		return ErrLongUserPassword
	}

	userRequest.Email = trimmedEmail
	userRequest.Password = trimmedPassword

	return nil
}
