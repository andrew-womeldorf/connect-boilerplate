package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	internalContext "github.com/andrew-womeldorf/connect-boilerplate/internal/context"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			requestID := req.Header().Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			ctx = internalContext.WithRequestID(ctx, requestID)

			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			res.Header().Set(RequestIDHeader, requestID)

			return res, nil
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
