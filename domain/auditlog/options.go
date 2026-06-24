package auditlog

// Option configures an API at construction time. Kept as a parameterless functional
// option set for forward compatibility — no API options ship in v1.
type Option[Payload any] func(*API[Payload])
