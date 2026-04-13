package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/MrBorisT/task-tracker-api/internal/auth"
	"github.com/MrBorisT/task-tracker-api/internal/helper"
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
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := verifyUserRequest(&userRequest); err != nil {
			log.Println("verifying user request:", err)
			_ = helper.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := userStore.RegisterUser(r.Context(), userRequest); err != nil {
			switch err {
			case storage.ErrUserAlreadyExists:
				_ = helper.WriteJSONError(w, http.StatusConflict, "user with this email already exists")
			default:
				log.Println("registering user:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func LoginUserHandler(userStore *storage.UserStore, authManager *auth.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRequest := models.UserRequest{}

		decoder := json.NewDecoder(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := decoder.Decode(&userRequest); err != nil {
			log.Println("decoding user:", err)
			_ = helper.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := verifyUserRequest(&userRequest); err != nil {
			log.Println("verifying user request:", err)
			_ = helper.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		userID, err := userStore.GetUserID(r.Context(), userRequest)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrUserNotFound):
				_ = helper.WriteJSONError(w, http.StatusNotFound, "user not found")
			case errors.Is(err, storage.ErrInvalidCredentials):
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "invalid email or password")
			default:
				log.Println("logging in user:", err)
				_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			}
			return
		}
		token, err := authManager.GenerateJWT(userID)
		if err != nil {
			log.Println("generate jwt:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(models.JWTToken{Token: token}); err != nil {
			log.Println("encoding jwt:", err)
			_ = helper.WriteJSONError(w, http.StatusInternalServerError, "something went wrong, try again later")
			return
		}
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
