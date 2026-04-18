package handlers

import (
	"testing"

	"github.com/MrBorisT/task-tracker-api/internal/models"
)

func TestVerifyUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *models.UserRequest
		err     error
	}{
		{
			name: "valid request",
			request: &models.UserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			err: nil,
		},
		{
			name: "empty email",
			request: &models.UserRequest{
				Email:    "",
				Password: "password123",
			},
			err: ErrEmptyUserEmail,
		},
		{
			name: "incorrect email",
			request: &models.UserRequest{
				Email:    "testexample.com",
				Password: "12345678",
			},
			err: ErrIncorrectUserEmail,
		},
		{
			name: "empty password",
			request: &models.UserRequest{
				Email:    "test@example.com",
				Password: "",
			},
			err: ErrEmptyUserPassword,
		},
		{
			name: "short password",
			request: &models.UserRequest{
				Email:    "test@example.com",
				Password: "123",
			},
			err: ErrShortUserPassword,
		},
		{
			name: "long password",
			request: &models.UserRequest{
				Email:    "test@example.com",
				Password: "9123456789123456789123456789123456789123456789123456789123456789123456789",
			},
			err: ErrLongUserPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyUserRequest(tt.request)

			if err != tt.err {
				t.Errorf("expected error: %v, got: %v", tt.err, err)
			}
		})
	}
}
