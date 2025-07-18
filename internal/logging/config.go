package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogformatter "github.com/samber/slog-formatter"
)

func SetupLogger(jsonFormat bool, level slog.Level) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if jsonFormat {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level: level,
		})
	}

	handler = slogformatter.NewFormatterHandler(
		slogformatter.ErrorFormatter("error"),
		slogformatter.PIIFormatter("email"),
	)(handler)

	handler = NewRequestIDHandler(handler)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
