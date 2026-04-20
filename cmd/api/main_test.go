package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MrBorisT/task-tracker-api/internal/auth"
	"github.com/MrBorisT/task-tracker-api/internal/config"
)

func TestGetTasks_WithoutToken(t *testing.T) {
	tests := []struct {
		name             string
		httpMethod       string
		requestTarget    string
		wantStatusCode   int
		wantBodyContains string
	}{
		{
			name:             "missing token",
			httpMethod:       http.MethodGet,
			requestTarget:    "/tasks",
			wantStatusCode:   http.StatusUnauthorized,
			wantBodyContains: "missing authorization header",
		},
		{
			name:             "invalid request body for login",
			httpMethod:       http.MethodPost,
			requestTarget:    "/auth/login",
			wantStatusCode:   http.StatusBadRequest,
			wantBodyContains: "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtManager := auth.NewJWTManager(&config.Config{
				JWTSecret: "test-secret",
				JWTTTL:    time.Hour,
			})

			router := newRouter(jwtManager, nil, nil)

			req := httptest.NewRequest(tt.httpMethod, tt.requestTarget, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			if !strings.Contains(rr.Body.String(), tt.wantBodyContains) {
				t.Fatalf("body = %q, want to contain %q", rr.Body.String(), tt.wantBodyContains)
			}
		})
	}
}
