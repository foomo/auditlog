// Package auditlogprivate carries the write-side wire types and the
// fully-formed generic Service shell mounted behind a project's
// non-generic internal Service interface.
//
// The package is named "private" rather than "internal" because Go's
// `internal/` directory rule would prevent the project module (different
// import-path root) from importing it. Semantics are unchanged: this
// surface is never mapped to TypeScript.
package auditlogprivate

// EntityID is the wire-form id used by Log. Kept as a string-typed alias so
// gotsrpc emits a plain TS `string` for the Go client.
type EntityID string

// DateTime is an ISO8601 millisecond-precision string. Duplicated from the
// public package so the two wire surfaces can drift independently.
type DateTime string

// Entry is the write-side envelope. Payload is the project-defined union.
// No bson tags — these never reach mongo; the repository works on
// store.Entry, and conv.go translates at the Service boundary.
type Entry[Payload any] struct {
	ID        EntityID `json:"id"`
	Timestamp DateTime `json:"timestamp"`
	UserID    string   `json:"userId"`
	RequestID string   `json:"requestId,omitempty"`
	Service   string   `json:"service"`
	Func      string   `json:"func"`
	Action    string   `json:"action"`
	EntityID  string   `json:"entityId,omitempty"`
	IP        string   `json:"ip,omitempty"`
	UserAgent string   `json:"userAgent,omitempty"`
	Payload   Payload  `json:"payload"`
}

// AuditLogError is the write-side error type returned across gotsrpc. Stringy
// so the Go client receives a plain string.
type AuditLogError string

// NewAuditLogError wraps a message as a pointer-returning constructor.
func NewAuditLogError(msg string) *AuditLogError {
	e := AuditLogError(msg)
	return &e
}

// Error satisfies the standard error interface.
func (e *AuditLogError) Error() string { return string(*e) }
