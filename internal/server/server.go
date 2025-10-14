package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/grpcreflect"
	sloghttp "github.com/samber/slog-http"
	slogmulti "github.com/samber/slog-multi"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	v1 "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1/userv1connect"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/web"
)

// Server represents the API server
type Server struct {
	port      int
	userStore store.Store
}

// NewServer creates a new server
func NewServer(port int, userStore store.Store) *Server {
	return &Server{
		port:      port,
		userStore: userStore,
	}
}

// Run starts the server
func (s *Server) Run() error {
	ctx := context.Background()

	configureLogging()

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

// Add common attributes to all logs
func configureLogging() {
	mid := slogmulti.NewHandleInlineMiddleware(func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
		if requestID := sloghttp.GetRequestIDFromContext(ctx); requestID != "" {
			record.AddAttrs(slog.String("id", requestID))
		}
		return next(ctx, record)
	})

	slog.SetDefault(slog.New(
		mid(slog.Default().Handler()),
	))
}

// CreateHandler creates an HTTP handler for the server without starting it
// This is useful for Lambda functions that need to handle HTTP requests
func (s *Server) CreateHandler(ctx context.Context) (http.Handler, error) {
	userService := user.NewService(s.userStore)
	webHandler := web.NewHandler(userService)

	// Create Connect server
	mux := http.NewServeMux()
	p, h := v1.NewUserServiceHandler(NewUserConnectHandler(userService))
	mux.Handle(p, h)

	// Add gRPC Reflector
	reflector := grpcreflect.NewStaticReflector(v1.UserServiceName)
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
	mux.HandleFunc("/", webHandler.IndexHandler)
	mux.HandleFunc("/create-user", webHandler.CreateUserHandler)

	// Add CORS middleware for browser clients
	mid := corsMiddleware(mux)
	mid = sloghttp.Recovery(mid)
	mid = sloghttp.NewWithConfig(
		slog.Default(),
		sloghttp.Config{
			WithSpanID:         true,
			WithTraceID:        true,
			WithRequestID:      true,
			WithRequestHeader:  true,
			WithResponseHeader: true,
		},
	)(mid)

	// Create h2c handler for HTTP/2 support
	h2cHandler := h2c.NewHandler(mid, &http2.Server{})

	return h2cHandler, nil
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
