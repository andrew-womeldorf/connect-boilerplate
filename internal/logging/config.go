package logging

import (
	"log/slog"
)

// SetupLogger adds a new handler to the default logger to add the request id
func SetupLogger() {
	handler := slog.Default().Handler()
	handler = NewRequestIDHandler(handler)
	slog.SetDefault(slog.New(handler))
}
