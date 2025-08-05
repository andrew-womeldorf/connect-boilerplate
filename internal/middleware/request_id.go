package middleware

import (
	"net/http"

	internalContext "github.com/andrew-womeldorf/connect-boilerplate/internal/context"
	sloghttp "github.com/samber/slog-http"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := internalContext.WithRequestID(
			r.Context(),
			sloghttp.GetRequestID(r),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
