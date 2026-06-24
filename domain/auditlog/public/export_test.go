package auditlogpublic

import (
	storex "github.com/foomo/auditlog/domain/auditlog/store"
)

// Test seams exposing the unexported conv helpers to the external
// auditlogpublic_test package.
func EntryFromStore[Payload any](e *storex.Entry[Payload]) *Entry[Payload] {
	return entryFromStore[Payload](e)
}

func PagedFromStore[Payload any](p *storex.PagedResult[storex.Entry[Payload]]) *PagedResult[Entry[Payload]] {
	return pagedFromStore[Payload](p)
}
