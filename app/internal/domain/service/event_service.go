package service

import (
	"awesomeProject1/internal/adapter"
	"awesomeProject1/internal/domain/model"
	"awesomeProject1/pkg/logger"
	"time"
)

type EventService struct {
	SNSPublisher adapter.SNSPublisher
}

func (s *EventService) ProcessAndPublishEvent(payload map[string]interface{}) error {
	log := logger.New()
	log.Debug("Starting to process event")

	event := model.Event{
		ID:          "unique-id",
		Payload:     payload,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	log.Infof("Event created with ID: %s", event.ID)

	err := s.SNSPublisher.Publish(event)
	if err != nil {
		log.Errorf("Failed to publish event: %v", err)
		return err
	}

	log.Info("Event published successfully")
	return nil
}
