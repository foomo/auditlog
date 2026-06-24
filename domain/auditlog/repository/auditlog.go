package auditlogrepository

import (
	"context"
	"time"

	storex "github.com/foomo/auditlog/domain/auditlog/store"
	keelmongo "github.com/foomo/keel/persistence/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

// AuditLogRepository is the persistence contract used by the command/query layer.
// It is generic over the project-defined payload type.
type AuditLogRepository[Payload any] interface {
	Insert(ctx context.Context, entry *storex.Entry[Payload]) error
	FindByID(ctx context.Context, id storex.EntityID) (*storex.Entry[Payload], error)
	Search(ctx context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error)
}

// BaseAuditLogRepository is the default mongo-backed implementation of
// AuditLogRepository.
type BaseAuditLogRepository[Payload any] struct {
	l          *zap.Logger
	collection *keelmongo.Collection
	retention  time.Duration
}

// NewBaseAuditLogRepository constructs a base repository, configuring the underlying
// mongo collection with a TTL index over the `ttlTime` field and a small set of
// compound indexes that match the supported Search filters.
func NewBaseAuditLogRepository[Payload any](
	l *zap.Logger,
	persistor *keelmongo.Persistor,
	opts ...Option,
) (*BaseAuditLogRepository[Payload], error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	collection, cErr := persistor.Collection(
		cfg.collectionName,
		keelmongo.CollectionWithIndexes(
			storex.EntityWithTTLIndex(cfg.retention, nil),
			mongo.IndexModel{
				Keys: bson.D{
					{Key: "service", Value: 1},
					{Key: "func", Value: 1},
					{Key: "timestamp", Value: -1},
				},
			},
			mongo.IndexModel{
				Keys: bson.D{
					{Key: "userId", Value: 1},
					{Key: "timestamp", Value: -1},
				},
			},
			mongo.IndexModel{
				Keys: bson.D{
					{Key: "entityId", Value: 1},
					{Key: "timestamp", Value: -1},
				},
			},
			mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		),
	)
	if cErr != nil {
		return nil, cErr
	}

	return &BaseAuditLogRepository[Payload]{
		l:          l,
		collection: collection,
		retention:  cfg.retention,
	}, nil
}

// Insert writes an entry to the collection, filling in default fields when not set
// by the caller (ID, Timestamp). TTLTime is always (re)anchored to "now" so that
// the TTL index measures retention from insert time.
func (r *BaseAuditLogRepository[Payload]) Insert(ctx context.Context, entry *storex.Entry[Payload]) error {
	if entry.ID == "" {
		entry.ID = storex.NewEntityID()
	}

	if entry.Timestamp == "" {
		entry.Timestamp = storex.NewDateTime(time.Now())
	}

	entry.TTLTime = time.Now()

	_, err := r.collection.Col().InsertOne(ctx, entry)

	return err
}

// FindByID returns the entry with the given id or mongo.ErrNoDocuments if none.
func (r *BaseAuditLogRepository[Payload]) FindByID(ctx context.Context, id storex.EntityID) (*storex.Entry[Payload], error) {
	var entry storex.Entry[Payload]

	if err := r.collection.FindOne(ctx, bson.M{"id": id}, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// Search runs a paginated, sorted query against the collection using the filter
// axes exposed by storex.Search. Empty / zero values are skipped.
func (r *BaseAuditLogRepository[Payload]) Search(ctx context.Context, qry storex.Search) (*storex.PagedResult[storex.Entry[Payload]], error) {
	page := max(qry.Page, 1)

	pageSize := qry.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	filter := bson.M{}

	if qry.Service != "" {
		filter["service"] = qry.Service
	}

	if qry.Func != "" {
		filter["func"] = qry.Func
	}

	if qry.Action != "" {
		filter["action"] = qry.Action
	}

	if qry.UserID != "" {
		filter["userId"] = qry.UserID
	}

	if qry.EntityID != "" {
		filter["entityId"] = qry.EntityID
	}

	if !qry.From.IsZero() || !qry.To.IsZero() {
		bounds := bson.M{}
		if !qry.From.IsZero() {
			bounds["$gte"] = storex.NewDateTime(qry.From)
		}

		if !qry.To.IsZero() {
			bounds["$lte"] = storex.NewDateTime(qry.To)
		}

		filter["timestamp"] = bounds
	}

	sortField := qry.Sort.Field
	if sortField == "" {
		sortField = storex.SortFieldTimestamp
	}

	skip := (page - 1) * pageSize
	findOpts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{
			{Key: string(sortField), Value: qry.Sort.Direction.GetSortValue()},
			{Key: "_id", Value: 1},
		})

	cursor, err := r.collection.Col().Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*storex.Entry[Payload]

	for cursor.Next(ctx) {
		var entry storex.Entry[Payload]
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}

		results = append(results, &entry)
	}

	total, err := r.collection.Col().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &storex.PagedResult[storex.Entry[Payload]]{
		Results:  results,
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}
