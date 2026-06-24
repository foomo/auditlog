package auditlogpublic_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auditlog "github.com/foomo/auditlog/domain/auditlog"
	auditlogpublic "github.com/foomo/auditlog/domain/auditlog/public"
	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testPayload struct {
	Source string `json:"source"`
}

// fakeRepo is an in-memory AuditLogRepository[testPayload] so the test can
// exercise Service end-to-end without mongo.
type fakeRepo struct {
	findResult *storex.Entry[testPayload]
	findErr    error
	searchRes  *storex.PagedResult[storex.Entry[testPayload]]
	searchErr  error
	lastSearch storex.Search
}

var _ repositoryx.AuditLogRepository[testPayload] = (*fakeRepo)(nil)

func (r *fakeRepo) Insert(_ context.Context, _ *storex.Entry[testPayload]) error {
	return errors.New("insert not used in public service tests")
}

func (r *fakeRepo) FindByID(_ context.Context, _ storex.EntityID) (*storex.Entry[testPayload], error) {
	return r.findResult, r.findErr
}

func (r *fakeRepo) Search(_ context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[testPayload]], error) {
	r.lastSearch = qry
	return r.searchRes, r.searchErr
}

func newService(t *testing.T, repo repositoryx.AuditLogRepository[testPayload]) *auditlogpublic.Service[testPayload] {
	t.Helper()

	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	return auditlogpublic.NewService[testPayload](zap.NewNop(), api)
}

func TestService_Get_returnsPublicEntry(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{findResult: &storex.Entry[testPayload]{
		ID:      "abc",
		Service: "redirects",
		Payload: testPayload{Source: "/old"},
	}}
	svc := newService(t, repo)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	entry, errResp := svc.Get(httptest.NewRecorder(), req, "abc")
	require.Nil(t, errResp)
	require.NotNil(t, entry)
	assert.Equal(t, auditlogpublic.EntityID("abc"), entry.ID)
	assert.Equal(t, "redirects", entry.Service)
	assert.Equal(t, testPayload{Source: "/old"}, entry.Payload)
}

func TestService_Get_translatesRepoError(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{findErr: errors.New("boom")}
	svc := newService(t, repo)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	entry, errResp := svc.Get(httptest.NewRecorder(), req, "abc")
	assert.Nil(t, entry)
	require.NotNil(t, errResp)
	assert.Equal(t, "boom", errResp.Error())
}

func TestService_Search_buildsStoreQueryAndConverts(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{searchRes: &storex.PagedResult[storex.Entry[testPayload]]{
		Results: []*storex.Entry[testPayload]{{ID: "x"}, {ID: "y"}},
		Total:   2, Page: 1, PageSize: 20,
	}}
	svc := newService(t, repo)

	params := &auditlogpublic.SearchParams{
		Service:  "redirects",
		Func:     "Delete",
		Action:   "delete",
		UserID:   "alice",
		EntityID: "redirect-1",
		Page:     1,
		PageSize: 20,
		Sort: auditlogpublic.Sort{
			Field:     auditlogpublic.SortFieldTimestamp,
			Direction: auditlogpublic.DirectionDescending,
		},
	}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	got, errResp := svc.Search(httptest.NewRecorder(), req, params)

	require.Nil(t, errResp)
	require.NotNil(t, got)
	assert.Equal(t, 2, got.Total)
	require.Len(t, got.Results, 2)
	assert.Equal(t, auditlogpublic.EntityID("x"), got.Results[0].ID)

	assert.Equal(t, "redirects", repo.lastSearch.Service)
	assert.Equal(t, "Delete", repo.lastSearch.Func)
	assert.Equal(t, "delete", repo.lastSearch.Action)
	assert.Equal(t, "alice", repo.lastSearch.UserID)
	assert.Equal(t, "redirect-1", repo.lastSearch.EntityID)
	assert.Equal(t, storex.SortFieldTimestamp, repo.lastSearch.Sort.Field)
	assert.Equal(t, storex.DirectionDescending, repo.lastSearch.Sort.Direction)
}

func TestService_Search_parsesFromAndTo(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{searchRes: &storex.PagedResult[storex.Entry[testPayload]]{}}
	svc := newService(t, repo)

	from := storex.NewDateTime(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	to := storex.NewDateTime(time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC))

	params := &auditlogpublic.SearchParams{
		From: auditlogpublic.DateTime(from),
		To:   auditlogpublic.DateTime(to),
	}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	_, errResp := svc.Search(httptest.NewRecorder(), req, params)
	require.Nil(t, errResp)

	assert.False(t, repo.lastSearch.From.IsZero())
	assert.False(t, repo.lastSearch.To.IsZero())
}

func TestService_Search_rejectsBadFrom(t *testing.T) {
	t.Parallel()

	svc := newService(t, &fakeRepo{})
	params := &auditlogpublic.SearchParams{From: "not-a-time"}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)

	got, errResp := svc.Search(httptest.NewRecorder(), req, params)
	assert.Nil(t, got)
	require.NotNil(t, errResp)
	assert.Contains(t, errResp.Error(), "invalid from")
}

func TestService_Search_rejectsBadTo(t *testing.T) {
	t.Parallel()

	svc := newService(t, &fakeRepo{})
	params := &auditlogpublic.SearchParams{To: "not-a-time"}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)

	got, errResp := svc.Search(httptest.NewRecorder(), req, params)
	assert.Nil(t, got)
	require.NotNil(t, errResp)
	assert.Contains(t, errResp.Error(), "invalid to")
}
