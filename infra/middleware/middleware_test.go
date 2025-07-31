package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"casino/utils"
)

type MockLogger struct{}

func (m *MockLogger) Error(ctx context.Context, errs ...error)     {}
func (m *MockLogger) Info(ctx context.Context, messages ...string) {}

func TestLoggingMiddleware(t *testing.T) {
	mockLogger := &MockLogger{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middlewareHandler := LoggingMiddleware(handler, mockLogger)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	if rr.Body.String() != "test response" {
		t.Errorf("Expected body 'test response', got %s", rr.Body.String())
	}
}

func TestLoggingMiddleware_WithRequestID(t *testing.T) {
	mockLogger := &MockLogger{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(utils.CtxKeyRequestID)
		if requestID == nil {
			t.Error("Expected request ID to be set")
		}
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := LoggingMiddleware(handler, mockLogger)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestLoggingMiddleware_ExistingRequestID(t *testing.T) {
	mockLogger := &MockLogger{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(utils.CtxKeyRequestID)
		if requestID != "existing-id" {
			t.Errorf("Expected request ID 'existing-id', got %v", requestID)
		}
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := LoggingMiddleware(handler, mockLogger)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(req.Context(), utils.CtxKeyRequestID, "existing-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}
