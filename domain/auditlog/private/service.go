package auditlogprivate

import (
	"net/http"

	auditlog "github.com/foomo/auditlog/domain/auditlog"
	"go.uber.org/zap"
)

// Service is the write-side gotsrpc handler. Generic over the project's
// Payload type so its instantiated method set matches the project's
// concrete internal Service interface — drop-in mountable.
type Service[Payload any] struct {
	l   *zap.Logger
	api *auditlog.API[Payload]
}

// NewService wires the write-side Service over an existing *API[Payload].
// One API instance can back both this Service and the read-side Service.
func NewService[Payload any](l *zap.Logger, api *auditlog.API[Payload], opts ...Option[Payload]) *Service[Payload] {
	s := &Service[Payload]{l: l, api: api}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Log writes a single audit entry. Auto-fill of ID / Timestamp / TTLTime
// happens inside the repository — callers may leave those fields zero.
func (s *Service[Payload]) Log(_ http.ResponseWriter, r *http.Request, entry *Entry[Payload]) *AuditLogError {
	if err := s.api.Log(r.Context(), entryToStore[Payload](entry)); err != nil {
		return NewAuditLogError(err.Error())
	}

	return nil
}
