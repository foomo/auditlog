# API reference

The authoritative API reference is hosted on pkg.go.dev:

- [`github.com/foomo/auditlog/domain/auditlog`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog) — `API[Payload]`, `NewAPI`, `Log` / `Get` / `Search`.
- [`.../domain/auditlog/store`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/store) — `Entry[Payload]`, `EntityID`, `DateTime`, `Search`, `Sort`, `PagedResult[T]`, `EntityWithTTL`.
- [`.../domain/auditlog/repository`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/repository) — `AuditLogRepository[Payload]`, `BaseAuditLogRepository[Payload]`, `WithRetention`, `WithCollectionName`.
- [`.../domain/auditlog/command`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/command) — `CreateEntry`, handler + middleware composer.
- [`.../domain/auditlog/query`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/query) — `Search`, `Get` handlers + middleware composers.
- [`.../domain/auditlog/public`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/public) — read-side `Service[Payload]` (`Get`, `Search`).
- [`.../domain/auditlog/private`](https://pkg.go.dev/github.com/foomo/auditlog/domain/auditlog/private) — write-side `Service[Payload]` (`Log`).

## `API[Payload]` methods

```go
// Log records an audit entry. The caller owns the entry fields except
// for ID / Timestamp / TTLTime, which the repository fills in when zero.
func (a *API[Payload]) Log(ctx context.Context, entry *storex.Entry[Payload]) error

// Get returns the entry with the given id.
func (a *API[Payload]) Get(ctx context.Context, id storex.EntityID) (*storex.Entry[Payload], error)

// Search returns a page of entries matching the filter.
func (a *API[Payload]) Search(ctx context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error)
```

For usage-oriented explanations with code examples, see the Guide:

- [Getting started](../guide/getting-started)
- [Architecture](../guide/architecture)
- [Retention](../guide/retention)
- [gotsrpc integration](../guide/gotsrpc)
