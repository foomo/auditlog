# Configuration reference

The repository is the only component with construction-time options.
`NewAPI` and the Services take a logger and the upstream dependency
positionally; everything tunable lives on the repository.

## `repository.NewBaseAuditLogRepository[Payload](l, persistor, opts ...Option)`

### Required parameters

| Parameter | Type | Purpose |
| --- | --- | --- |
| `l` | `*zap.Logger` | Logger handed to the repository. |
| `persistor` | `*keelmongo.Persistor` | MongoDB persistor that owns the connection and collection factory. |

### Options

| Option | Default | Purpose |
| --- | --- | --- |
| `WithRetention(d time.Duration)` | `180 * 24h` | How long entries are kept before the MongoDB TTL index removes them. Baked into the TTL index — see [Retention](../guide/retention). |
| `WithCollectionName(name string)` | `auditlog` | Overrides the MongoDB collection name. |

### Indexes created

`NewBaseAuditLogRepository` configures the collection with:

- a **TTL index** over `ttlTime` (`expireAfterSeconds` from `WithRetention`);
- a **unique index** over `id`;
- compound indexes matching the `Search` filter axes:
  - `{ service: 1, func: 1, timestamp: -1 }`
  - `{ userId: 1, timestamp: -1 }`
  - `{ entityId: 1, timestamp: -1 }`

## `auditlog.NewAPI[Payload](l, repo, opts ...Option[Payload])`

| Parameter | Type | Purpose |
| --- | --- | --- |
| `l` | `*zap.Logger` | Logger. Returns an error if nil. |
| `repo` | `AuditLogRepository[Payload]` | Persistence backend. Returns an error if nil. |
| `opts` | `Option[Payload]` | Functional options. **None ship in v1** — the type exists for forward compatibility. |

`NewAPI` applies no middleware. Layer telemetry / event publishing /
capability checks via the `command` and `query` middleware composers.

## `public.NewService` / `private.NewService`

| Parameter | Type | Purpose |
| --- | --- | --- |
| `l` | `*zap.Logger` | Logger. |
| `api` | `*auditlog.API[Payload]` | The shared API instance. One API can back both Services. |
| `opts` | `Option[Payload]` | Functional options for the Service. |

## `store.Search`

The filter object accepted by `API.Search`. All fields are optional;
empty / zero values are skipped.

| Field | Type | Effect |
| --- | --- | --- |
| `Service` | `string` | Exact match on `service`. |
| `Func` | `string` | Exact match on `func`. |
| `Action` | `string` | Exact match on `action`. |
| `UserID` | `string` | Exact match on `userId`. |
| `EntityID` | `string` | Exact match on `entityId`. |
| `From` | `time.Time` | `timestamp >= From`. |
| `To` | `time.Time` | `timestamp <= To`. |
| `Page` | `int` | 1-based page; defaults to `1` when `< 1`. |
| `PageSize` | `int` | Page size; defaults to `20` when `< 1`. |
| `Sort` | `Sort` | Sort field + direction; defaults to `timestamp`. |
