package bootstrap

import (
	"awesomeProject1/internal/adapter"
	httpHandler "awesomeProject1/internal/adapter/http"
	"awesomeProject1/internal/core/domain/service"
	"awesomeProject1/pkg/cfg"
	"awesomeProject1/pkg/logger"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaApplication struct {
	config       *cfg.Config
	awsService   *cfg.AWSService
	eventService *service.EventService
	eventHandler *httpHandler.EventHandler
	logger       *logger.Logger
}

func NewLambdaApplication() (*LambdaApplication, error) {
	log := logger.New()
	log.Info("Bootstrapping Lambda application")

	config, err := cfg.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	log.Infof("Application: %s, Environment: %s", config.App.Name, config.App.Environment)

	awsService, err := cfg.NewAWSService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AWS services: %w", err)
	}

	if config.IsLocalStack() {
		log.Infof("Using LocalStack endpoint: %s", config.AWS.LocalStack.Endpoint)
	} else {
		log.Info("Using AWS SNS service")
	}

	snsPublisher := &adapter.SNSPublisher{
		Client:   awsService.GetSNSClient(),
		TopicArn: awsService.GetSNSTopicArn(),
	}

	log.Infof("SNS Topic ARN: %s", awsService.GetSNSTopicArn())

	eventService := service.NewEventService(snsPublisher)
	eventHandler := httpHandler.NewEventHandler(eventService, log)

	app := &LambdaApplication{
		config:       config,
		awsService:   awsService,
		eventService: eventService,
		eventHandler: eventHandler,
		logger:       log,
	}

	log.Info("Lambda application bootstrapped successfully")
	return app, nil
}

func (app *LambdaApplication) Handler() func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return app.eventHandler.HandleLambda
}
