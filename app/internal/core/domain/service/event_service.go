package service

import (
	"awesomeProject1/internal/core/domain/model"
	"awesomeProject1/internal/core/port"
	"awesomeProject1/pkg/logger"
	"encoding/json"
	"fmt"
	"time"
)

type EventService struct {
	Publisher port.EventOutboundPort
}

func NewEventService(publisher port.EventOutboundPort) *EventService {
	return &EventService{
		Publisher: publisher,
	}
}

func (s *EventService) ProcessEventosOperacao(eventosOperacao *model.EventosOperacao) error {
	log := logger.New()
	log.Debugf("Starting to process eventos operacao of type: %s", eventosOperacao.EventType)

	if !eventosOperacao.IsValid() {
		log.Error("Invalid eventos operacao: missing required person identifiers")
		return fmt.Errorf("invalid eventos operacao: at least one of idPessoaFisica or idPessoaJuridica must be provided")
	}

	eventosOperacao.ProcessedAt = time.Now()
	eventosOperacao.EventType = eventosOperacao.DetermineEventType()

	log.Infof("Processing eventos operacao of type: %s for metadata.id: %s",
		eventosOperacao.EventType, eventosOperacao.Metadata.Id)

	payloadMap, err := s.convertToMap(eventosOperacao)
	if err != nil {
		log.Errorf("Failed to convert domain model to map: %v", err)
		return err
	}

	err = s.Publisher.PublishEvent(eventosOperacao.EventType, payloadMap)
	if err != nil {
		log.Errorf("Failed to publish event: %v", err)
		return err
	}

	log.Info("Eventos operacao processed and published successfully")
	return nil
}

func (s *EventService) convertToMap(eventosOperacao *model.EventosOperacao) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(eventosOperacao)
	if err != nil {
		return nil, err
	}

	var payloadMap map[string]interface{}
	err = json.Unmarshal(jsonData, &payloadMap)
	if err != nil {
		return nil, err
	}

	return payloadMap, nil
}
