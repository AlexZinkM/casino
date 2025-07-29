package logging

import (
	"casino/boundary/logging"
	"casino/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

type logType string

const (
	logTypeInfo  logType = "log"
	logTypeError logType = "error"
)

type logMessage struct {
	ctx     context.Context
	logType logType
	errs    []error
	msgs    []string
}

type AsyncLogger struct {
	appName string
	ch      chan logMessage
	mu      sync.RWMutex
	loggers []logging.Logger
}

func NewAsyncLogger(appName string) *AsyncLogger {
	f := &AsyncLogger{
		ch:      make(chan logMessage, 2),
		appName: appName,
	}
	go f.run()
	return f
}

func (f *AsyncLogger) Register(l logging.Logger) {
	f.mu.Lock()
	f.loggers = append(f.loggers, l)
	f.mu.Unlock()
}

func (f *AsyncLogger) Error(ctx context.Context, errs ...error) {
	f.ch <- logMessage{ctx: ctx, logType: logTypeError, errs: errs}
}

func (f *AsyncLogger) Info(ctx context.Context, msgs ...string) {
	f.ch <- logMessage{ctx: ctx, logType: logTypeInfo, msgs: msgs}
}

func (f *AsyncLogger) Close() {
	close(f.ch)
}

func (f *AsyncLogger) run() {
	for m := range f.ch {
		f.mu.RLock()
		for _, l := range f.loggers {
			if m.logType == logTypeError {
				l.Error(withAppName(m.ctx, f.appName), m.errs...)
			} else {
				l.Info(withAppName(m.ctx, f.appName), m.msgs...)
			}
		}
		f.mu.RUnlock()
	}
}

type SimpleLogger struct{}

func (s *SimpleLogger) Error(ctx context.Context, errs ...error) {
	fmt.Printf("[%s] ERROR [%s] [%s] %v\n", appNameFromCtx(ctx), requestIDFromCtx(ctx), time.Now().Format(time.RFC3339), errs)
}

func (s *SimpleLogger) Info(ctx context.Context, msgs ...string) {
	fmt.Printf("[%s] INFO  [%s] [%s] %v\n", appNameFromCtx(ctx), requestIDFromCtx(ctx), time.Now().Format(time.RFC3339), msgs)
}

type ctxKey string

const ctxKeyAppName ctxKey = "appName"

func withAppName(ctx context.Context, appName string) context.Context {
	return context.WithValue(ctx, ctxKeyAppName, appName)
}

func appNameFromCtx(ctx context.Context) string {
	v := ctx.Value(ctxKeyAppName)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func requestIDFromCtx(ctx context.Context) string {
	v := ctx.Value(utils.CtxKeyRequestID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
