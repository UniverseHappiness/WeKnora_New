package test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
)

func newCtxWithBuffer() (context.Context, *bytes.Buffer) {
	ctx := context.Background()
	entry := logger.GetLogger(ctx)
	buf := &bytes.Buffer{}
	entry.Logger.SetOutput(buf)
	ctx = context.WithValue(ctx, types.LoggerContextKey, entry)
	return ctx, buf
}

func TestLoggerInfoIncludesCaller(t *testing.T) {
	ctx, buf := newCtxWithBuffer()
	ctx = logger.WithRequestID(ctx, "req-123")
	logger.Info(ctx, "hello world")

	out := buf.String()
	if !strings.Contains(out, "hello world") {
		t.Fatalf("log output missing message: %s", out)
	}
	if !strings.Contains(out, "[TestLoggerInfoIncludesCaller]") {
		t.Fatalf("log output missing caller function name: %s", out)
	}
	if !strings.Contains(out, "request_id=req-123") {
		t.Fatalf("log output missing request_id field: %s", out)
	}
}

func TestLoggerInfofFormatting(t *testing.T) {
	ctx, buf := newCtxWithBuffer()
	ctx = logger.WithRequestID(ctx, "req-xyz")
	logger.Infof(ctx, "value=%d %s", 42, "ok")

	out := buf.String()
	if !strings.Contains(out, "value=42 ok") {
		t.Fatalf("log output missing formatted message: %s", out)
	}
	if !strings.Contains(out, "request_id=req-xyz") {
		t.Fatalf("log output missing request_id field: %s", out)
	}
}

func TestLoggerFieldsOrdering(t *testing.T) {
	ctx, buf := newCtxWithBuffer()
	ctx = logger.WithRequestID(ctx, "r1")
	ctx = logger.WithField(ctx, "extra", 1)
	logger.Info(ctx, "x")

	out := buf.String()
	firstOpen := strings.Index(out, "[")
	if firstOpen == -1 {
		t.Fatalf("missing first '[' in output: %s", out)
	}
	firstClose := strings.Index(out[firstOpen+1:], "]")
	if firstClose == -1 {
		t.Fatalf("missing first ']' in output: %s", out)
	}
	secondOpen := strings.Index(out[firstOpen+firstClose+2:], "[")
	if secondOpen == -1 {
		t.Fatalf("missing second '[' (fields) in output: %s", out)
	}
	secondOpen += firstOpen + firstClose + 2
	secondClose := strings.Index(out[secondOpen+1:], "]")
	if secondClose == -1 {
		t.Fatalf("missing second ']' (fields) in output: %s", out)
	}
	fields := out[secondOpen+1 : secondOpen+1+secondClose]
	if !strings.HasPrefix(fields, "request_id=r1") {
		t.Fatalf("request_id not leading fields: %s", fields)
	}
	if !strings.Contains(fields, "extra=1") {
		t.Fatalf("missing extra field: %s", fields)
	}
}
