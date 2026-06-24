package auditlogstore

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// EntityWithTTLIndex returns a mongo index model that expires documents based on the
// `ttlTime` field after the given duration. Mirrors the commerce primitive so the
// library does not depend on the commerce module.
func EntityWithTTLIndex(ttl time.Duration, opts *options.IndexOptionsBuilder) mongo.IndexModel {
	expireAfterSeconds := int32(ttl.Seconds())

	if opts == nil {
		opts = options.Index()
	}

	opts = opts.SetExpireAfterSeconds(expireAfterSeconds)

	return mongo.IndexModel{
		Keys: bson.D{
			{Key: "ttlTime", Value: 1},
		},
		Options: opts,
	}
}

// EntityWithTTL embeds the TTL anchor field used by the TTL index.
type EntityWithTTL struct {
	TTLTime time.Time `json:"ttlTime" bson:"ttlTime" yaml:"ttlTime"`
}

// EnsureTTLTime sets TTLTime to now when zero.
func (e *EntityWithTTL) EnsureTTLTime() {
	if e.TTLTime.IsZero() {
		e.TTLTime = time.Now()
	}
}
