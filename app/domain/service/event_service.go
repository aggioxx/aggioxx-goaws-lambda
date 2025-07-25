package service

import (
	"awesomeProject1/adapter"
	"awesomeProject1/domain/model"
	"time"
)

type EventService struct {
	SNSPublisher adapter.SNSPublisher
}

func (s *EventService) ProcessAndPublishEvent(payload map[string]interface{}) error {
	event := model.Event{
		ID:          "unique-id",
		Payload:     payload,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	return s.SNSPublisher.Publish(event)
}
