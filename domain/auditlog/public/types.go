// Package auditlogpublic carries the read-side wire types and the
// fully-formed generic Service shell that mounts behind a project's
// non-generic public Service interface.
//
// These types never reach mongo; the repository layer works on the
// separate store/ package. Public and private wire-side packages are
// deliberately duplicated so the two surfaces can drift independently.
package auditlogpublic

import "time"

// EntityID is the wire-form id used by Get/Search. Kept as a string-typed
// alias so gotsrpc emits a plain TS `string`.
type EntityID string

// DateTime is an ISO8601 millisecond-precision string. Same semantics as
// store.DateTime; duplicated here to keep the public package free of
// storage-layer dependencies.
type DateTime string

// dateTimeLayout matches storex.DateTimeLayout so the two string-typed
// DateTime values are wire-compatible.
const dateTimeLayout = "2006-01-02T15:04:05.000Z0700"

// Time parses the DateTime in ISO8601 ms format.
func (d DateTime) Time() (time.Time, error) {
	return time.Parse(dateTimeLayout, string(d))
}

// SortField enumerates the entry fields the read endpoint can sort by.
type SortField string

const (
	SortFieldTimestamp SortField = "timestamp"
	SortFieldService   SortField = "service"
	SortFieldFunc      SortField = "func"
	SortFieldAction    SortField = "action"
	SortFieldUserID    SortField = "userId"
)

// Direction is the sort order over a SortField.
type Direction string

const (
	DirectionAscending  Direction = "ascending"
	DirectionDescending Direction = "descending"
)

// Sort is the sort instruction carried inside SearchParams.
type Sort struct {
	Field     SortField `json:"field"`
	Direction Direction `json:"direction"`
}

// Entry is the read-side envelope. Payload is the project-defined union.
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

// SearchParams is the over-the-wire input shape for Service.Search. From
// and To are DateTime strings so the gotsrpc-generated TS client surfaces
// them as plain strings.
type SearchParams struct {
	Service  string   `json:"service,omitempty"`
	Func     string   `json:"func,omitempty"`
	Action   string   `json:"action,omitempty"`
	UserID   string   `json:"userId,omitempty"`
	EntityID string   `json:"entityId,omitempty"`
	From     DateTime `json:"from,omitempty"`
	To       DateTime `json:"to,omitempty"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Sort     Sort     `json:"sort"`
}

// PagedResult is the paginated response wrapper for Search.
type PagedResult[T any] struct {
	Results  []*T `json:"results"`
	Total    int  `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
}

// AuditLogError is the read-side error type returned across gotsrpc. Stringy
// so the TS client receives a plain string.
type AuditLogError string

// NewAuditLogError wraps a message as a pointer-returning constructor in the
// same shape gotsrpc consumers expect.
func NewAuditLogError(msg string) *AuditLogError {
	e := AuditLogError(msg)
	return &e
}

// Error satisfies the standard error interface.
func (e *AuditLogError) Error() string { return string(*e) }
