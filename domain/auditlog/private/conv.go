package auditlogprivate

import (
	storex "github.com/foomo/auditlog/domain/auditlog/store"
)

// entryToStore maps a write-side wire entry to its repository-layer form.
// The store.Entry's embedded EntityWithTTL is left zero — the repository
// (Insert) anchors TTLTime to time.Now() before insert.
func entryToStore[Payload any](e *Entry[Payload]) *storex.Entry[Payload] {
	if e == nil {
		return nil
	}

	return &storex.Entry[Payload]{
		ID:        storex.EntityID(e.ID),
		Timestamp: storex.DateTime(e.Timestamp),
		UserID:    e.UserID,
		RequestID: e.RequestID,
		Service:   e.Service,
		Func:      e.Func,
		Action:    e.Action,
		EntityID:  e.EntityID,
		IP:        e.IP,
		UserAgent: e.UserAgent,
		Payload:   e.Payload,
	}
}
