package logging

import (
	"context"
	"fmt"
	"log/slog"

	internalContext "github.com/andrew-womeldorf/connect-boilerplate/internal/context"
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
	fmt.Println("before handle logger")
	if requestID, ok := internalContext.GetRequestID(ctx); ok {
		record.AddAttrs(slog.String("id", requestID))
	}
	fmt.Println("after handle logger")
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
