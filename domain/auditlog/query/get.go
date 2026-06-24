package auditlogquery

import (
	"context"

	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	// GetHandlerFn handler
	GetHandlerFn[Payload any] func(ctx context.Context, l *zap.Logger, id storex.EntityID) (*storex.Entry[Payload], error)
	// GetMiddlewareFn middleware
	GetMiddlewareFn[Payload any] func(next GetHandlerFn[Payload]) GetHandlerFn[Payload]
)

// GetHandler is the leaf handler — it just delegates to the repository.
func GetHandler[Payload any](repo repositoryx.AuditLogRepository[Payload]) GetHandlerFn[Payload] {
	return func(ctx context.Context, _ *zap.Logger, id storex.EntityID) (*storex.Entry[Payload], error) {
		return repo.FindByID(ctx, id)
	}
}

// GetHandlerComposed returns the leaf handler wrapped in the provided middlewares.
func GetHandlerComposed[Payload any](
	handler GetHandlerFn[Payload],
	middlewares ...GetMiddlewareFn[Payload],
) GetHandlerFn[Payload] {
	composed := func(next GetHandlerFn[Payload]) GetHandlerFn[Payload] {
		for _, middleware := range middlewares {
			localNext := next
			middlewareName := funcName(middleware)
			next = middleware(func(ctx context.Context, l *zap.Logger, id storex.EntityID) (*storex.Entry[Payload], error) {
				trace.SpanFromContext(ctx).AddEvent(middlewareName)
				return localNext(ctx, l, id)
			})
		}

		return next
	}
	handlerName := funcName(handler)

	return composed(func(ctx context.Context, l *zap.Logger, id storex.EntityID) (*storex.Entry[Payload], error) {
		trace.SpanFromContext(ctx).AddEvent(handlerName)
		return handler(ctx, l, id)
	})
}
