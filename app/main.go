package main

import (
	"awesomeProject1/adapter"
	"awesomeProject1/domain/service"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
)

var eventService *service.EventService

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	snsClient := sns.NewFromConfig(cfg)
	snsPublisher := adapter.SNSPublisher{
		Client:   snsClient,
		TopicArn: "arn:aws:sns:us-east-1:000000000000:my-topic",
	}
	eventService = &service.EventService{
		SNSPublisher: snsPublisher,
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" && request.Path == "/callbacks-trackingagrovarejo/v1/eventos_operacao" {
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(request.Body), &payload)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "` + err.Error() + `"}`,
			}, nil
		}

		err = eventService.ProcessAndPublishEvent(payload)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "` + err.Error() + `"}`,
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"message": "Event processed and published successfully"}`,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error": "Not Found"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
