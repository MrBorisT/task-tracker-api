package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MrBorisT/task-tracker-api/internal/auth"
	"github.com/MrBorisT/task-tracker-api/internal/config"
)

func TestAuthMiddleware(t *testing.T) {
	jwtManager := auth.NewJWTManager(&config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    time.Hour,
	})

	validToken, err := jwtManager.GenerateJWT("user-123")
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	tests := []struct {
		name              string
		authHeader        string
		wantStatus        int
		wantBodyContains  string
		wantNextCalled    bool
		wantContextUserID string
	}{
		{
			name:             "missing authorization header",
			authHeader:       "",
			wantStatus:       http.StatusUnauthorized,
			wantBodyContains: "missing authorization header",
			wantNextCalled:   false,
		},
		{
			name:             "invalid authorization header prefix",
			authHeader:       "Basic abc123",
			wantStatus:       http.StatusUnauthorized,
			wantBodyContains: "invalid authorization header",
			wantNextCalled:   false,
		},
		{
			name:             "missing token after bearer",
			authHeader:       "Bearer ",
			wantStatus:       http.StatusUnauthorized,
			wantBodyContains: "missing token",
			wantNextCalled:   false,
		},
		{
			name:              "valid token",
			authHeader:        "Bearer " + validToken,
			wantStatus:        http.StatusOK,
			wantNextCalled:    true,
			wantContextUserID: "user-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nextCalled bool
			var gotUserID string

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				userID, ok := r.Context().Value(UserIDKey).(string)
				if ok {
					gotUserID = userID
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok":true}`))
			})

			handler := AuthMiddleware(jwtManager)(next)

			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rr.Code, tt.wantStatus, rr.Body.String())
			}

			if nextCalled != tt.wantNextCalled {
				t.Fatalf("nextCalled = %v, want %v", nextCalled, tt.wantNextCalled)
			}

			if tt.wantBodyContains != "" && !strings.Contains(rr.Body.String(), tt.wantBodyContains) {
				t.Fatalf("body = %q, want to contain %q", rr.Body.String(), tt.wantBodyContains)
			}

			if gotUserID != tt.wantContextUserID {
				t.Fatalf("context userID = %q, want %q", gotUserID, tt.wantContextUserID)
			}
		})
	}
}
