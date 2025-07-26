package main

import (
	"awesomeProject1/pkg/bootstrap"
	"awesomeProject1/pkg/logger"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func main() {
	log := logger.New()

	app, err := bootstrap.NewLambdaApplication()
	if err != nil {
		log.Errorf("Failed to bootstrap application: %v", err)
		os.Exit(1)
	}

	log.Info("Starting Lambda handler")
	lambda.Start(app.Handler())
}
