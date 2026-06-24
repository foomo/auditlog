package auditlogprivate

import (
	storex "github.com/foomo/auditlog/domain/auditlog/store"
)

// Test seam exposing the unexported conv helper to the external
// auditlogprivate_test package.
func EntryToStore[Payload any](e *Entry[Payload]) *storex.Entry[Payload] {
	return entryToStore[Payload](e)
}
