package logging

import (
	"context"
	"log/slog"

	internal_context "github.com/andrew-womeldorf/connect-boilerplate/internal/context"
)

type RequestIDHandler struct {
	handler slog.Handler
}

func NewRequestIDHandler(handler slog.Handler) *RequestIDHandler {
	return &RequestIDHandler{
		handler: handler,
	}
}

func (h *RequestIDHandler) Handle(ctx context.Context, record slog.Record) error {
	if requestID, ok := internal_context.GetRequestID(ctx); ok {
		record.AddAttrs(slog.String("id", requestID))
	}
	return h.handler.Handle(ctx, record)
}

func (h *RequestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewRequestIDHandler(h.handler.WithAttrs(attrs))
}

func (h *RequestIDHandler) WithGroup(name string) slog.Handler {
	return NewRequestIDHandler(h.handler.WithGroup(name))
}

func (h *RequestIDHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}
