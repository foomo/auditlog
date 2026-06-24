package auditlogpublic

import (
	"net/http"

	auditlog "github.com/foomo/auditlog/domain/auditlog"
	storex "github.com/foomo/auditlog/domain/auditlog/store"
	"go.uber.org/zap"
)

// Service is the read-side gotsrpc handler. It is generic over the
// project's Payload type; method receivers' signatures, once instantiated,
// match the project's concrete Service interface — so the project mounts
// *Service[Payload] directly without writing a wrapper struct.
type Service[Payload any] struct {
	l   *zap.Logger
	api *auditlog.API[Payload]
}

// NewService wires the read-side Service over an existing *API[Payload].
// One API instance can back both this Service and the write-side Service.
func NewService[Payload any](l *zap.Logger, api *auditlog.API[Payload], opts ...Option[Payload]) *Service[Payload] {
	s := &Service[Payload]{l: l, api: api}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Get returns the entry with the given id.
func (s *Service[Payload]) Get(_ http.ResponseWriter, r *http.Request, id string) (*Entry[Payload], *AuditLogError) {
	entry, err := s.api.Get(r.Context(), storex.EntityID(id))
	if err != nil {
		return nil, NewAuditLogError(err.Error())
	}

	return entryFromStore[Payload](entry), nil
}

// Search returns a paged, filtered, sorted slice of entries.
func (s *Service[Payload]) Search(_ http.ResponseWriter, r *http.Request, params *SearchParams) (*PagedResult[Entry[Payload]], *AuditLogError) {
	qry := storex.Search{
		Service:  params.Service,
		Func:     params.Func,
		Action:   params.Action,
		UserID:   params.UserID,
		EntityID: params.EntityID,
		Page:     params.Page,
		PageSize: params.PageSize,
		Sort: storex.Sort{
			Field:     storex.SortField(params.Sort.Field),
			Direction: storex.Direction(params.Sort.Direction),
		},
	}
	if params.From != "" {
		t, err := params.From.Time()
		if err != nil {
			return nil, NewAuditLogError("invalid from: " + err.Error())
		}

		qry.From = t
	}

	if params.To != "" {
		t, err := params.To.Time()
		if err != nil {
			return nil, NewAuditLogError("invalid to: " + err.Error())
		}

		qry.To = t
	}

	page, err := s.api.Search(r.Context(), qry)
	if err != nil {
		return nil, NewAuditLogError(err.Error())
	}

	return pagedFromStore[Payload](page), nil
}
