package main

// Logs mit log/slog: Text oder JSON (NOTES_LOG_FORMAT=json), dazu trace_id und
// span_id, sobald eine Anfrage einen Trace hat (Tag 11).

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func newLogger() *slog.Logger {
	var h slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if os.Getenv("NOTES_LOG_FORMAT") == "json" {
		h = slog.NewJSONHandler(os.Stdout, nil)
	}
	return slog.New(traceHandler{h})
}

type traceHandler struct{ slog.Handler }

func (h traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()))
	}
	return h.Handler.Handle(ctx, r)
}

func (h traceHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return traceHandler{h.Handler.WithAttrs(as)}
}

func (h traceHandler) WithGroup(name string) slog.Handler {
	return traceHandler{h.Handler.WithGroup(name)}
}

func logRequest(r *http.Request, code int, took time.Duration) {
	level := slog.LevelInfo
	if code >= 500 {
		level = slog.LevelError
	}
	slog.Log(r.Context(), level, "request", "method", r.Method, "path", r.URL.Path,
		"status", code, "duration_ms", took.Milliseconds())
}
