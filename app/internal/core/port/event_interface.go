package port

import "awesomeProject1/internal/core/domain/model"

type EventInboundPort interface {
	ProcessEventosOperacao(eventosOperacao *model.EventosOperacao) error
}

type EventOutboundPort interface {
	PublishEvent(eventType string, payload map[string]interface{}) error
}
