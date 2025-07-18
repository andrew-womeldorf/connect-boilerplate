package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/andrew-womeldorf/connect-boilerplate/gen/example/v1/examplev1connect"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/interceptor"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/web"
	"github.com/andrew-womeldorf/connect-boilerplate/pkg/api"
)

// Server represents the API server
type Server struct {
	port int
}

// NewServer creates a new server
func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

// Run starts the server
func (s *Server) Run() error {
	ctx := context.Background()

	// Create handler
	handler, err := s.CreateHandler(ctx)
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}

	// Start server
	addr := fmt.Sprintf(":%d", s.port)
	slog.InfoContext(ctx, "server listening", slog.String("address", addr))
	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// CreateHandler creates an HTTP handler for the server without starting it
// This is useful for Lambda functions that need to handle HTTP requests
func (s *Server) CreateHandler(ctx context.Context) (http.Handler, error) {
	// Create service
	serviceHandler, err := s.createService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	// Get the underlying service for web handlers
	service, err := s.GetService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Create Connect server
	mux := http.NewServeMux()
	path, connectHandler := examplev1connect.NewUserServiceHandler(serviceHandler,
		connect.WithInterceptors(
			interceptor.RequestIDInterceptor(),
		),
	)
	mux.Handle(path, connectHandler)

	// Add gRPC Reflector
	reflector := grpcreflect.NewStaticReflector(examplev1connect.UserServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.ErrorContext(ctx, "Failed to write health check response", slog.Any("error", err))
		}
	})

	// Add web interface endpoints
	webHandler := web.NewHandler(service)
	mux.HandleFunc("/", webHandler.IndexHandler)
	mux.HandleFunc("/create-user", webHandler.CreateUserHandler)

	// Add CORS middleware for browser clients
	corsHandler := corsMiddleware(mux)

	// Create h2c handler for HTTP/2 support
	h2cHandler := h2c.NewHandler(corsHandler, &http2.Server{})

	return h2cHandler, nil
}

// GetService returns the service from the server
func (s *Server) GetService(ctx context.Context) (*api.Service, error) {
	serviceHandler, err := s.createService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	// Get the service adapter
	serviceAdapter, ok := serviceHandler.(*ServiceAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to cast service to ServiceAdapter")
	}

	return serviceAdapter.service, nil
}

// createService creates the service and adapter
func (s *Server) createService(ctx context.Context) (examplev1connect.UserServiceHandler, error) {
	// Create service
	service := api.NewService()

	// Create adapter
	adapter := NewServiceAdapter(service)

	return adapter, nil
}

// corsMiddleware adds CORS headers for browser clients
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
