package auditlog_test

import (
	"context"
	"errors"
	"testing"

	auditlog "github.com/foomo/auditlog/domain/auditlog"
	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// testPayload is a project-shaped union payload for the API tests.
type testPayload struct {
	Redirect *testRedirectPayload `json:"redirect,omitempty"`
}

type testRedirectPayload struct {
	Source string `json:"source"`
}

// mockRepo records calls and returns canned responses. It implements
// repositoryx.AuditLogRepository[testPayload].
type mockRepo struct {
	inserted   []*storex.Entry[testPayload]
	insertErr  error
	findResult *storex.Entry[testPayload]
	findErr    error
	searchRes  *storex.PagedResult[storex.Entry[testPayload]]
	searchErr  error
	lastSearch storex.Search
}

var _ repositoryx.AuditLogRepository[testPayload] = (*mockRepo)(nil)

func (m *mockRepo) Insert(_ context.Context, e *storex.Entry[testPayload]) error {
	m.inserted = append(m.inserted, e)
	return m.insertErr
}
func (m *mockRepo) FindByID(_ context.Context, _ storex.EntityID) (*storex.Entry[testPayload], error) {
	return m.findResult, m.findErr
}
func (m *mockRepo) Search(_ context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[testPayload]], error) {
	m.lastSearch = qry
	return m.searchRes, m.searchErr
}

func TestNewAPI_validatesArgs(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		l    *zap.Logger
		repo repositoryx.AuditLogRepository[testPayload]
	}{
		"missing logger":     {l: nil, repo: &mockRepo{}},
		"missing repository": {l: zap.NewNop(), repo: nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := auditlog.NewAPI[testPayload](tc.l, tc.repo)
			assert.Error(t, err)
		})
	}
}

func TestAPI_Log_delegatesToRepo(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{}
	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	entry := &storex.Entry[testPayload]{
		Service: "redirects",
		Func:    "Delete",
		Action:  "delete",
		Payload: testPayload{Redirect: &testRedirectPayload{Source: "/a"}},
	}

	require.NoError(t, api.Log(context.Background(), entry))
	require.Len(t, repo.inserted, 1)
	assert.Same(t, entry, repo.inserted[0])
}

func TestAPI_Log_returnsRepoError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	repo := &mockRepo{insertErr: wantErr}
	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	gotErr := api.Log(context.Background(), &storex.Entry[testPayload]{})
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAPI_Get_delegatesToRepo(t *testing.T) {
	t.Parallel()

	want := &storex.Entry[testPayload]{ID: "abc"}
	repo := &mockRepo{findResult: want}
	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	got, err := api.Get(context.Background(), "abc")
	require.NoError(t, err)
	assert.Same(t, want, got)
}

func TestAPI_Search_passesQuery(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{searchRes: &storex.PagedResult[storex.Entry[testPayload]]{}}
	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	qry := storex.Search{Service: "redirects", PageSize: 5}
	_, err = api.Search(context.Background(), qry)
	require.NoError(t, err)
	assert.Equal(t, qry, repo.lastSearch)
}
