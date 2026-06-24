package auditlogpublic

// Option configures a Service at construction time. Reserved for forward
// compatibility — no options ship in v1.
type Option[Payload any] func(*Service[Payload])
