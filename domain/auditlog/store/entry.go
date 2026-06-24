package auditlogstore

// Entry is the audit log envelope. Payload is project-defined and typically a union
// struct of pointers to per-entity payload variants.
type Entry[Payload any] struct {
	ID            EntityID `json:"id" bson:"id"`
	Timestamp     DateTime `json:"timestamp" bson:"timestamp"`
	UserID        string   `json:"userId" bson:"userId"`
	RequestID     string   `json:"requestId,omitempty" bson:"requestId,omitempty"`
	Service       string   `json:"service" bson:"service"`
	Func          string   `json:"func" bson:"func"`
	Action        string   `json:"action" bson:"action"`
	EntityID      string   `json:"entityId,omitempty" bson:"entityId,omitempty"`
	IP            string   `json:"ip,omitempty" bson:"ip,omitempty"`
	UserAgent     string   `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	EntityWithTTL `json:",inline" bson:",inline"`
	Payload       Payload `json:"payload" bson:"payload"`
}
