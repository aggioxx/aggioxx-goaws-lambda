package main

import (
	"awesomeProject1/internal/adapter"
	"awesomeProject1/internal/domain/service"
	"awesomeProject1/pkg/logger"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var eventService *service.EventService

func init() {
	log := logger.New()
	log.Info("Initializing application")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Errorf("Failed to load configuration: %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)
	snsPublisher := adapter.SNSPublisher{
		Client:   snsClient,
		TopicArn: "arn:aws:sns:us-east-1:000000000000:my-topic",
	}
	eventService = &service.EventService{
		SNSPublisher: snsPublisher,
	}

	log.Info("Application initialized successfully")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log := logger.New()
	log.Debugf("Received request: %v", request)

	if request.HTTPMethod == "POST" && request.Path == "/callbacks-trackingagrovarejo/v1/eventos_operacao" {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(request.Body), &payload)
		if err != nil {
			log.Warnf("Invalid request body: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "` + err.Error() + `"}`,
			}, nil
		}

		log.Info("Processing event")
		err = eventService.ProcessAndPublishEvent(payload)
		if err != nil {
			log.Errorf("Failed to process and publish event: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "` + err.Error() + `"}`,
			}, nil
		}

		log.Info("Event processed and published successfully")
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Event processed and published successfully"}`,
		}, nil
	}

	log.Warn("Unhandled request path")
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error": "Not Found"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
