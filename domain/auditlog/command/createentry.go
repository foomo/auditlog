package auditlogcommand

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
	// CreateEntry command
	CreateEntry[Payload any] struct {
		Entry *storex.Entry[Payload] `json:"entry"`
	}
	// CreateEntryHandlerFn handler
	CreateEntryHandlerFn[Payload any] func(ctx context.Context, l *zap.Logger, cmd CreateEntry[Payload]) error
	// CreateEntryMiddlewareFn middleware
	CreateEntryMiddlewareFn[Payload any] func(next CreateEntryHandlerFn[Payload]) CreateEntryHandlerFn[Payload]
)

// CreateEntryHandler is the leaf handler — it just delegates to the repository.
func CreateEntryHandler[Payload any](repo repositoryx.AuditLogRepository[Payload]) CreateEntryHandlerFn[Payload] {
	return func(ctx context.Context, _ *zap.Logger, cmd CreateEntry[Payload]) error {
		return repo.Insert(ctx, cmd.Entry)
	}
}

// CreateEntryHandlerComposed returns the leaf handler wrapped in the provided
// middlewares; each middleware call emits a span event for tracing.
func CreateEntryHandlerComposed[Payload any](
	handler CreateEntryHandlerFn[Payload],
	middlewares ...CreateEntryMiddlewareFn[Payload],
) CreateEntryHandlerFn[Payload] {
	composed := func(next CreateEntryHandlerFn[Payload]) CreateEntryHandlerFn[Payload] {
		for _, middleware := range middlewares {
			localNext := next
			middlewareName := funcName(middleware)
			next = middleware(func(ctx context.Context, l *zap.Logger, cmd CreateEntry[Payload]) error {
				trace.SpanFromContext(ctx).AddEvent(middlewareName)
				return localNext(ctx, l, cmd)
			})
		}

		return next
	}
	handlerName := funcName(handler)

	return composed(func(ctx context.Context, l *zap.Logger, cmd CreateEntry[Payload]) error {
		trace.SpanFromContext(ctx).AddEvent(handlerName)
		return handler(ctx, l, cmd)
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
