package auditlogprivate_test

import (
	"testing"

	auditlogprivate "github.com/foomo/auditlog/domain/auditlog/private"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryToStore_copiesAllFields(t *testing.T) {
	t.Parallel()

	in := &auditlogprivate.Entry[testPayload]{
		ID:        "id-1",
		Timestamp: "2026-05-22T10:00:00.000Z",
		UserID:    "alice",
		RequestID: "req-1",
		Service:   "redirects",
		Func:      "Delete",
		Action:    "delete",
		EntityID:  "redirect-1",
		IP:        "10.0.0.1",
		UserAgent: "Mozilla/5.0",
		Payload:   testPayload{Source: "/old"},
	}

	got := auditlogprivate.EntryToStore[testPayload](in)

	require.NotNil(t, got)
	assert.Equal(t, storex.EntityID("id-1"), got.ID)
	assert.Equal(t, storex.DateTime("2026-05-22T10:00:00.000Z"), got.Timestamp)
	assert.Equal(t, "alice", got.UserID)
	assert.Equal(t, "req-1", got.RequestID)
	assert.Equal(t, "redirects", got.Service)
	assert.Equal(t, "Delete", got.Func)
	assert.Equal(t, "delete", got.Action)
	assert.Equal(t, "redirect-1", got.EntityID)
	assert.Equal(t, "10.0.0.1", got.IP)
	assert.Equal(t, "Mozilla/5.0", got.UserAgent)
	assert.Equal(t, testPayload{Source: "/old"}, got.Payload)
}

func TestEntryToStore_returnsNilOnNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, auditlogprivate.EntryToStore[testPayload](nil))
}
