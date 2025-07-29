package logging

import (
	"context"
	"fmt"
	"testing"

	"casino/utils"
)

func TestNewAsyncLogger(t *testing.T) {
	logger := NewAsyncLogger("test-app")

	if logger == nil {
		t.Error("Expected logger to be created")
	}

	if logger.appName != "test-app" {
		t.Errorf("Expected appName 'test-app', got %s", logger.appName)
	}

	if logger.ch == nil {
		t.Error("Expected channel to be created")
	}

	defer logger.Close()
}

func TestAsyncLogger_Register(t *testing.T) {
	logger := NewAsyncLogger("test-app")
	simpleLogger := &SimpleLogger{}

	logger.Register(simpleLogger)

	if len(logger.loggers) != 1 {
		t.Errorf("Expected 1 logger, got %d", len(logger.loggers))
	}

	defer logger.Close()
}

func TestAsyncLogger_Error(t *testing.T) {
	logger := NewAsyncLogger("test-app")
	simpleLogger := &SimpleLogger{}
	logger.Register(simpleLogger)

	ctx := context.Background()
	err := fmt.Errorf("test error")

	logger.Error(ctx, err)

	defer logger.Close()
}

func TestAsyncLogger_Info(t *testing.T) {
	logger := NewAsyncLogger("test-app")
	simpleLogger := &SimpleLogger{}
	logger.Register(simpleLogger)

	ctx := context.Background()
	message := "test message"

	logger.Info(ctx, message)

	defer logger.Close()
}

func TestSimpleLogger_Error(t *testing.T) {
	logger := &SimpleLogger{}
	ctx := context.Background()
	err := fmt.Errorf("test error")

	logger.Error(ctx, err)
}

func TestSimpleLogger_Info(t *testing.T) {
	logger := &SimpleLogger{}
	ctx := context.Background()
	message := "test message"

	logger.Info(ctx, message)
}

func TestWithAppName(t *testing.T) {
	ctx := context.Background()
	appName := "test-app"

	newCtx := withAppName(ctx, appName)

	if newCtx == ctx {
		t.Error("Expected new context to be created")
	}
}

func TestAppNameFromCtx(t *testing.T) {
	ctx := context.Background()
	appName := "test-app"

	newCtx := withAppName(ctx, appName)
	retrievedAppName := appNameFromCtx(newCtx)

	if retrievedAppName != appName {
		t.Errorf("Expected appName %s, got %s", appName, retrievedAppName)
	}
}

func TestRequestIDFromCtx(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	newCtx := context.WithValue(ctx, utils.CtxKeyRequestID, requestID)
	retrievedRequestID := requestIDFromCtx(newCtx)

	if retrievedRequestID != requestID {
		t.Errorf("Expected requestID %s, got %s", requestID, retrievedRequestID)
	}
}
