# Getting started

`auditlog` gives you a generic, type-safe audit-log domain: a fixed
`Entry[Payload]` envelope around a payload union *you* define, persisted
in MongoDB with automatic TTL-based retention. The library follows the
same shape as `foomo/redirects` — a generic `*API[Payload]` composes
command and query handlers over a generic `Repository[Payload]`.

## Install

```sh
go get github.com/foomo/auditlog
```

The pieces you wire together live under `domain/auditlog`:

```go
import (
    auditlog "github.com/foomo/auditlog/domain/auditlog"
    auditlogprivate "github.com/foomo/auditlog/domain/auditlog/private"
    auditlogpublic "github.com/foomo/auditlog/domain/auditlog/public"
    auditrepo "github.com/foomo/auditlog/domain/auditlog/repository"
    storex "github.com/foomo/auditlog/domain/auditlog/store"
)
```

## 1. Declare your payload

The payload is project-defined and typically a union of pointers to
per-entity variants. Declare it once, next to the rest of your audit-log
service code.

```go
package myaudit

import (
    redirectstore "github.com/foomo/redirects/v2/domain/redirectdefinition/store"
)

type AuditLog struct {
    Redirect *AuditLogRedirect `json:"redirect,omitempty"`
}

type AuditLogRedirect struct {
    Before *redirectstore.RedirectDefinition `json:"before,omitempty"`
    After  *redirectstore.RedirectDefinition `json:"after,omitempty"`
}
```

## 2. Wire up the API

One `*API[AuditLog]` backs both the read-side and write-side Services.

```go
ctx := context.Background()
l, _ := zap.NewProduction()

persistor, err := cmrcmongo.New(ctx, "mongodb://localhost:27017/local")
log.Must(l, err, "failed to create persistor")

repo, err := auditrepo.NewBaseAuditLogRepository[myaudit.AuditLog](l, persistor,
    auditrepo.WithRetention(180*24*time.Hour),
)
log.Must(l, err, "failed to create audit log repository")

api, err := auditlog.NewAPI[myaudit.AuditLog](l, repo)
log.Must(l, err, "failed to create audit log api")

publicService := auditlogpublic.NewService[myaudit.AuditLog](l, api)
privateService := auditlogprivate.NewService[myaudit.AuditLog](l, api)
```

## 3. Log, get, search

In-process callers use the `*API` directly:

```go
err := api.Log(ctx, &storex.Entry[myaudit.AuditLog]{
    Service: "redirects",
    Func:    "DeleteRedirect",
    Action:  "delete",
    UserID:  "alice",
    Payload: myaudit.AuditLog{Redirect: &myaudit.AuditLogRedirect{ /* ... */ }},
})
```

`ID`, `Timestamp` and `TTLTime` are filled in by the repository when
left zero — see [Retention](./retention).

## Where to go next

- [Architecture](./architecture) — the generic layering, what each package owns, and the invariants that hold it together.
- [Retention](./retention) — how the TTL index works and the one caveat when you change the horizon.
- [gotsrpc integration](./gotsrpc) — the public/private Service split and why the library ships no wire surface.
- [Configuration reference](../reference/configuration) — every repository option in one place.
