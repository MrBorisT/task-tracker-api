package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MrBorisT/task-tracker-api/internal/auth"
	"github.com/MrBorisT/task-tracker-api/internal/config"
	"github.com/MrBorisT/task-tracker-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func TestGetTasks_WithoutToken(t *testing.T) {
	jwtManager := auth.NewJWTManager(&config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    time.Hour,
	})

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(jwtManager))

	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tasks":[]}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusUnauthorized, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "missing authorization header") {
		t.Fatalf("body = %q, want to contain %q", rr.Body.String(), "missing authorization header")
	}
}
