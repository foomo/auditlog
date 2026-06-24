package auditlogpublic

import (
	storex "github.com/foomo/auditlog/domain/auditlog/store"
)

// entryFromStore maps a repository-layer entry to its read-side wire form.
// The embedded EntityWithTTL on storex.Entry is intentionally dropped — TTL
// is a storage concern that the public surface does not expose.
func entryFromStore[Payload any](e *storex.Entry[Payload]) *Entry[Payload] {
	if e == nil {
		return nil
	}

	return &Entry[Payload]{
		ID:        EntityID(e.ID),
		Timestamp: DateTime(e.Timestamp),
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

func pagedFromStore[Payload any](p *storex.PagedResult[storex.Entry[Payload]]) *PagedResult[Entry[Payload]] {
	if p == nil {
		return nil
	}

	out := &PagedResult[Entry[Payload]]{
		Total:    p.Total,
		Page:     p.Page,
		PageSize: p.PageSize,
		Results:  make([]*Entry[Payload], 0, len(p.Results)),
	}
	for _, e := range p.Results {
		out.Results = append(out.Results, entryFromStore[Payload](e))
	}

	return out
}
