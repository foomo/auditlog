package auditlog

import (
	"context"
	"errors"

	commandx "github.com/foomo/auditlog/domain/auditlog/command"
	queryx "github.com/foomo/auditlog/domain/auditlog/query"
	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"go.uber.org/zap"
)

// API is the domain entry point. It composes the command/query handlers around a
// project-typed repository and exposes the methods used by both the in-process
// caller and the *Service wrapper that backs gotsrpc.
type API[Payload any] struct {
	l    *zap.Logger
	qry  Queries[Payload]
	cmd  Commands[Payload]
	repo repositoryx.AuditLogRepository[Payload]
}

// NewAPI wires the command and query pipelines against the provided repository.
// No middlewares are applied in the library — telemetry and event publishing are
// project concerns.
func NewAPI[Payload any](
	l *zap.Logger,
	repo repositoryx.AuditLogRepository[Payload],
	opts ...Option[Payload],
) (*API[Payload], error) {
	if l == nil {
		return nil, errors.New("missing logger")
	}

	if repo == nil {
		return nil, errors.New("missing repository")
	}

	inst := &API[Payload]{
		l:    l,
		repo: repo,
	}

	for _, opt := range opts {
		opt(inst)
	}

	inst.cmd = Commands[Payload]{
		CreateEntry: commandx.CreateEntryHandlerComposed[Payload](
			commandx.CreateEntryHandler[Payload](inst.repo),
		),
	}
	inst.qry = Queries[Payload]{
		Search: queryx.SearchHandlerComposed[Payload](
			queryx.SearchHandler[Payload](inst.repo),
		),
		Get: queryx.GetHandlerComposed[Payload](
			queryx.GetHandler[Payload](inst.repo),
		),
	}

	return inst, nil
}

// Log records an audit entry. The caller owns the entry fields except for ID /
// Timestamp / TTLTime which the repository fills in when zero.
func (a *API[Payload]) Log(ctx context.Context, entry *storex.Entry[Payload]) error {
	return a.cmd.CreateEntry(ctx, a.l, commandx.CreateEntry[Payload]{Entry: entry})
}

// Get returns the entry with the given id.
func (a *API[Payload]) Get(ctx context.Context, id storex.EntityID) (*storex.Entry[Payload], error) {
	return a.qry.Get(ctx, a.l, id)
}

// Search returns a page of entries matching the filter.
func (a *API[Payload]) Search(ctx context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error) {
	return a.qry.Search(ctx, a.l, qry)
}
