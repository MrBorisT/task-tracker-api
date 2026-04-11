package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrEmptyUserEmail     = errors.New("user email cannot be empty")
	ErrIncorrectUserEmail = errors.New("user email must be a valid email address")
	ErrEmptyUserPassword  = errors.New("user password cannot be empty")
	ErrShortUserPassword  = errors.New("user password must be at least 6 characters long")
	ErrLongUserPassword   = errors.New("user password cannot be longer than 72 characters")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
