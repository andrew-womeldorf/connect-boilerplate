package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/andrew-womeldorf/connect-boilerplate/internal/server"
)

func init() {
	// Configure logger
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// Create server
	srv := server.NewServer(0) // Port doesn't matter for lambda

	// Create handler
	ctx := context.Background()
	handler, err := srv.CreateHandler(ctx)
	if err != nil {
		slog.Error("Failed to create handler", "error", err)
		os.Exit(1)
	}

	// Register the Lambda handler
	lambda.Start(handler)
}
