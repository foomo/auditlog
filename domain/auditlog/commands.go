package auditlog

import (
	commandx "github.com/foomo/auditlog/domain/auditlog/command"
)

type Commands[Payload any] struct {
	CreateEntry commandx.CreateEntryHandlerFn[Payload]
}
