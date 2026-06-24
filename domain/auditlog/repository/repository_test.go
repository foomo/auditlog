package auditlogrepository_test

// Integration test scaffold for the audit log repository. Mirrors the shape of
// foomo/redirects' repository_test.go — kept commented because CI here does not
// run a mongo instance. Uncomment and adjust the connection URI to run locally.

// import (
// 	"context"
// 	"testing"
// 	"time"
//
// 	auditlogrepository "github.com/foomo/auditlog/domain/auditlog/repository"
// 	auditlogstore "github.com/foomo/auditlog/domain/auditlog/store"
// 	keelmongo "github.com/foomo/keel/persistence/mongo"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
// 	"go.uber.org/zap"
// )
//
// type testPayload struct {
// 	Kind string `json:"kind" bson:"kind"`
// 	Note string `json:"note" bson:"note"`
// }
//
// func newTestRepo(t *testing.T) *auditlogrepository.BaseAuditLogRepository[testPayload] {
// 	t.Helper()
// 	l := zap.NewNop()
// 	persistor, err := keelmongo.New(
// 		context.Background(),
// 		"mongodb://localhost:27017/auditlog_test",
// 		keelmongo.WithOtelOptions(otelmongo.WithCommandAttributeDisabled(true)),
// 	)
// 	require.NoError(t, err)
// 	repo, err := auditlogrepository.NewBaseAuditLogRepository[testPayload](
// 		l, persistor,
// 		auditlogrepository.WithRetention(24*time.Hour),
// 		auditlogrepository.WithCollectionName("auditlog_test"),
// 	)
// 	require.NoError(t, err)
// 	return repo
// }
//
// func TestInsertAndFindByID(t *testing.T) {
// 	repo := newTestRepo(t)
// 	ctx := context.Background()
// 	entry := &auditlogstore.Entry[testPayload]{
// 		Service: "redirects",
// 		Func:    "DeleteRedirect",
// 		Action:  "delete",
// 		UserID:  "alice",
// 		Payload: testPayload{Kind: "redirect", Note: "demo"},
// 	}
// 	require.NoError(t, repo.Insert(ctx, entry))
// 	assert.NotEmpty(t, entry.ID)
// 	assert.NotEmpty(t, entry.Timestamp)
// 	assert.False(t, entry.TTLTime.IsZero())
//
// 	got, err := repo.FindByID(ctx, entry.ID)
// 	require.NoError(t, err)
// 	assert.Equal(t, entry.UserID, got.UserID)
// 	assert.Equal(t, "demo", got.Payload.Note)
// }
//
// func TestSearchByFilters(t *testing.T) {
// 	repo := newTestRepo(t)
// 	ctx := context.Background()
//
// 	// Insert a small fixture set.
// 	require.NoError(t, repo.Insert(ctx, &auditlogstore.Entry[testPayload]{Service: "redirects", Func: "Delete", Action: "delete", UserID: "alice", EntityID: "r-1"}))
// 	require.NoError(t, repo.Insert(ctx, &auditlogstore.Entry[testPayload]{Service: "redirects", Func: "Create", Action: "create", UserID: "alice", EntityID: "r-2"}))
// 	require.NoError(t, repo.Insert(ctx, &auditlogstore.Entry[testPayload]{Service: "customer", Func: "Update", Action: "update", UserID: "bob", EntityID: "c-1"}))
//
// 	for name, qry := range map[string]auditlogstore.Search{
// 		"service":  {Service: "redirects"},
// 		"func":     {Func: "Create"},
// 		"userId":   {UserID: "bob"},
// 		"entityId": {EntityID: "r-1"},
// 	} {
// 		t.Run(name, func(t *testing.T) {
// 			res, err := repo.Search(ctx, qry)
// 			require.NoError(t, err)
// 			assert.NotEmpty(t, res.Results)
// 		})
// 	}
// }
