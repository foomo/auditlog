package auditlogprivate_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auditlog "github.com/foomo/auditlog/domain/auditlog"
	auditlogprivate "github.com/foomo/auditlog/domain/auditlog/private"
	repositoryx "github.com/foomo/auditlog/domain/auditlog/repository"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testPayload struct {
	Source string `json:"source"`
}

type fakeRepo struct {
	inserted  []*storex.Entry[testPayload]
	insertErr error
}

var _ repositoryx.AuditLogRepository[testPayload] = (*fakeRepo)(nil)

func (r *fakeRepo) Insert(_ context.Context, e *storex.Entry[testPayload]) error {
	r.inserted = append(r.inserted, e)
	return r.insertErr
}
func (r *fakeRepo) FindByID(_ context.Context, _ storex.EntityID) (*storex.Entry[testPayload], error) {
	return nil, errors.New("not used")
}
func (r *fakeRepo) Search(_ context.Context, _ storex.Search) (*storex.PagedResult[storex.Entry[testPayload]], error) {
	return nil, errors.New("not used")
}

func newService(t *testing.T, repo repositoryx.AuditLogRepository[testPayload]) *auditlogprivate.Service[testPayload] {
	t.Helper()

	api, err := auditlog.NewAPI[testPayload](zap.NewNop(), repo)
	require.NoError(t, err)

	return auditlogprivate.NewService[testPayload](zap.NewNop(), api)
}

func TestService_Log_insertsConvertedEntry(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{}
	svc := newService(t, repo)

	entry := &auditlogprivate.Entry[testPayload]{
		ID:      "id-1",
		Service: "redirects",
		Func:    "Delete",
		Action:  "delete",
		Payload: testPayload{Source: "/old"},
	}
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)

	errResp := svc.Log(httptest.NewRecorder(), req, entry)
	require.Nil(t, errResp)
	require.Len(t, repo.inserted, 1)
	assert.Equal(t, storex.EntityID("id-1"), repo.inserted[0].ID)
	assert.Equal(t, "redirects", repo.inserted[0].Service)
	assert.Equal(t, testPayload{Source: "/old"}, repo.inserted[0].Payload)
}

func TestService_Log_translatesRepoError(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{insertErr: errors.New("boom")}
	svc := newService(t, repo)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	errResp := svc.Log(httptest.NewRecorder(), req, &auditlogprivate.Entry[testPayload]{})
	require.NotNil(t, errResp)
	assert.Equal(t, "boom", errResp.Error())
}
