package auditlogrepository

import "time"

// DefaultRetention is used when no WithRetention option is provided.
const DefaultRetention = 180 * 24 * time.Hour

// DefaultCollectionName is the mongo collection name used when no WithCollectionName
// option is provided.
const DefaultCollectionName = "auditlog"

type config struct {
	retention      time.Duration
	collectionName string
}

func defaultConfig() config {
	return config{
		retention:      DefaultRetention,
		collectionName: DefaultCollectionName,
	}
}

// Option configures a base audit log repository at construction time.
type Option func(*config)

// WithRetention sets how long audit entries are kept before the mongo TTL index
// deletes them. The duration is converted to whole seconds for the TTL index.
func WithRetention(d time.Duration) Option {
	return func(c *config) {
		c.retention = d
	}
}

// WithCollectionName overrides the mongo collection name (defaults to "auditlog").
func WithCollectionName(name string) Option {
	return func(c *config) {
		c.collectionName = name
	}
}
