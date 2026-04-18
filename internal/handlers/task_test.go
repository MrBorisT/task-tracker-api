package handlers

import (
	"testing"

	"github.com/MrBorisT/task-tracker-api/internal/models"
)

func TestNewGetTasksQuery(t *testing.T) {
	tests := []struct {
		name      string
		statusStr string
		limitStr  string
		want      *models.GetTasksQuery
		wantErr   bool
	}{
		{
			name:      "valid status and limit",
			statusStr: "new",
			limitStr:  "5",
			want: &models.GetTasksQuery{
				Status: "new",
				Limit:  5,
			},
			wantErr: false,
		},
		{
			name:      "empty limit uses default",
			statusStr: "done",
			limitStr:  "",
			want: &models.GetTasksQuery{
				Status: "done",
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name:      "invalid status",
			statusStr: "abc",
			limitStr:  "5",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "invalid limit string",
			statusStr: "new",
			limitStr:  "lol",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "zero limit",
			statusStr: "new",
			limitStr:  "0",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "negative limit",
			statusStr: "new",
			limitStr:  "-1",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newGetTasksQuery(tt.statusStr, tt.limitStr)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Status != tt.want.Status {
				t.Errorf("Status = %q, want %q", got.Status, tt.want.Status)
			}
			if got.Limit != tt.want.Limit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.want.Limit)
			}
		})
	}
}
