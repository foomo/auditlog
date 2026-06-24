package auditlogpublic_test

import (
	"testing"

	auditlogpublic "github.com/foomo/auditlog/domain/auditlog/public"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryFromStore_copiesAllFields(t *testing.T) {
	t.Parallel()

	in := &storex.Entry[testPayload]{
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

	got := auditlogpublic.EntryFromStore[testPayload](in)

	require.NotNil(t, got)
	assert.Equal(t, auditlogpublic.EntityID("id-1"), got.ID)
	assert.Equal(t, auditlogpublic.DateTime("2026-05-22T10:00:00.000Z"), got.Timestamp)
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

func TestEntryFromStore_returnsNilOnNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, auditlogpublic.EntryFromStore[testPayload](nil))
}

func TestPagedFromStore_copiesEntriesAndMetadata(t *testing.T) {
	t.Parallel()

	in := &storex.PagedResult[storex.Entry[testPayload]]{
		Results: []*storex.Entry[testPayload]{
			{ID: "a"},
			{ID: "b"},
		},
		Total:    42,
		Page:     2,
		PageSize: 10,
	}

	got := auditlogpublic.PagedFromStore[testPayload](in)

	require.NotNil(t, got)
	assert.Equal(t, 42, got.Total)
	assert.Equal(t, 2, got.Page)
	assert.Equal(t, 10, got.PageSize)
	require.Len(t, got.Results, 2)
	assert.Equal(t, auditlogpublic.EntityID("a"), got.Results[0].ID)
	assert.Equal(t, auditlogpublic.EntityID("b"), got.Results[1].ID)
}

func TestPagedFromStore_returnsNilOnNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, auditlogpublic.PagedFromStore[testPayload](nil))
}
