package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Kabanya/YAFDS/pkg/models"
)

func TestNumThreads(t *testing.T) {
	maxProcs := runtime.GOMAXPROCS(0)
	tests := []struct {
		name      string
		requested int
		want      uint8
	}{
		{
			name:      "requested zero",
			requested: 0,
			want:      uint8(maxProcs),
		},
		{
			name:      "requested negative",
			requested: -1,
			want:      uint8(maxProcs),
		},
		{
			name:      "requested one",
			requested: 1,
			want:      1,
		},
		{
			name:      "requested more than max",
			requested: maxProcs + 1,
			want:      uint8(maxProcs),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumThreads(tt.requested); got != tt.want {
				t.Errorf("NumThreads(%d) = %v, want %v", tt.requested, got, tt.want)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	content := `
DB_HOST=localhost
DB_PORT=5432
APP_NAME=YAFDS
EXPANDED=$(APP_NAME)_ENV
# Comment line
EMPTY_VAL=
`
	tmpFile := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp .env file: %v", err)
	}

	err := LoadEnv(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnv failed: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"APP_NAME", "YAFDS"},
		{"EXPANDED", "YAFDS_ENV"},
		{"EMPTY_VAL", ""},
	}

	for _, tt := range tests {
		got := os.Getenv(tt.key)
		if got != tt.want {
			t.Errorf("Env %s = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestLogger(t *testing.T) {
	tmpLog := filepath.Join(t.TempDir(), "test.log")

	// Test Init
	err := InitFileLogger(tmpLog)
	if err != nil {
		t.Fatalf("InitFileLogger failed: %v", err)
	}

	// Test Logger retrieval
	l, err := Logger()
	if err != nil {
		t.Fatalf("Logger() failed: %v", err)
	}
	if l == nil {
		t.Fatal("Logger() returned nil")
	}

	// Test writing
	l.Println("test log message")

	// Verify file exists and has content
	content, err := os.ReadFile(tmpLog)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if len(content) == 0 {
		t.Error("log file is empty")
	}

	// Test Close
	err = CloseLogger()
	if err != nil {
		t.Errorf("CloseLogger failed: %v", err)
	}
}

func TestUUID(t *testing.T) {
	// Test NewUUID
	id := NewUUID()
	if id == UuidNil {
		t.Error("NewUUID returned nil UUID")
	}

	// Test ParseUUID
	idStr := id.String()
	parsed, err := ParseUUID(idStr)
	if err != nil {
		t.Errorf("ParseUUID failed: %v", err)
	}
	if parsed != id {
		t.Errorf("ParseUUID = %v, want %v", parsed, id)
	}

	// Test ParseUUID invalid
	_, err = ParseUUID("invalid-uuid")
	if err == nil {
		t.Error("ParseUUID should fail for invalid string")
	}
}

func TestHTTPResponses(t *testing.T) {
	t.Run("WriteJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "hello"}
		WriteJSON(w, data, http.StatusOK)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("content-type = %q, want %q", contentType, "application/json")
		}

		var got map[string]string
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got["message"] != "hello" {
			t.Errorf("got message %q, want %q", got["message"], "hello")
		}
	})

	t.Run("WriteError", func(t *testing.T) {
		w := httptest.NewRecorder()
		errMsg := "some error"
		WriteError(w, errMsg, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var got models.ErrorResponce
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ErrorMessage != errMsg {
			t.Errorf("got error %q, want %q", got.ErrorMessage, errMsg)
		}
	})
}
