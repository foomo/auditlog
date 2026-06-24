package auditlogquery

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	// SearchHandlerFn handler
	SearchHandlerFn[Payload any] func(ctx context.Context, l *zap.Logger, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error)
	// SearchMiddlewareFn middleware
	SearchMiddlewareFn[Payload any] func(next SearchHandlerFn[Payload]) SearchHandlerFn[Payload]
)

// SearchHandler is the leaf handler — it just delegates to the repository.
func SearchHandler[Payload any](repo repositoryx.AuditLogRepository[Payload]) SearchHandlerFn[Payload] {
	return func(ctx context.Context, _ *zap.Logger, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error) {
		return repo.Search(ctx, qry)
	}
}

// SearchHandlerComposed returns the leaf handler wrapped in the provided middlewares.
func SearchHandlerComposed[Payload any](
	handler SearchHandlerFn[Payload],
	middlewares ...SearchMiddlewareFn[Payload],
) SearchHandlerFn[Payload] {
	composed := func(next SearchHandlerFn[Payload]) SearchHandlerFn[Payload] {
		for _, middleware := range middlewares {
			localNext := next
			middlewareName := funcName(middleware)
			next = middleware(func(ctx context.Context, l *zap.Logger, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error) {
				trace.SpanFromContext(ctx).AddEvent(middlewareName)
				return localNext(ctx, l, qry)
			})
		}

		return next
	}
	handlerName := funcName(handler)

	return composed(func(ctx context.Context, l *zap.Logger, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error) {
		trace.SpanFromContext(ctx).AddEvent(handlerName)
		return handler(ctx, l, qry)
	})
}

func funcName(fn any) string {
	full := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	parts := strings.Split(full, ".")
	if len(parts) < 3 {
		return full
	}

	return parts[2]
}
