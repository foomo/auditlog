package auditlogstore

import "time"

// Search captures all filter axes supported by the repository's Search.
// Empty string / zero time fields are treated as "no filter".
type Search struct {
	Service  string    `json:"service,omitempty"`
	Func     string    `json:"func,omitempty"`
	Action   string    `json:"action,omitempty"`
	UserID   string    `json:"userId,omitempty"`
	EntityID string    `json:"entityId,omitempty"`
	From     time.Time `json:"from,omitzero"`
	To       time.Time `json:"to,omitzero"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Sort     Sort      `json:"sort"`
}
