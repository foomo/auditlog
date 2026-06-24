package auditlog

import (
	queryx "github.com/foomo/auditlog/domain/auditlog/query"
)

type Queries[Payload any] struct {
	Search queryx.SearchHandlerFn[Payload]
	Get    queryx.GetHandlerFn[Payload]
}
