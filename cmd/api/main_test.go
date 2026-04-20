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
	jwtManager := auth.NewJWTManager(&config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    time.Hour,
	})

	router := newRouter(jwtManager, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusUnauthorized, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "missing authorization header") {
		t.Fatalf("body = %q, want to contain %q", rr.Body.String(), "missing authorization header")
	}
}
