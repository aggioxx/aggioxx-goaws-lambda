package adapter

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSPublisher struct {
	Client   *sns.Client
	TopicArn string
}

func (p *SNSPublisher) PublishEvent(eventType string, payload map[string]interface{}) error {
	message, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = p.Client.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(p.TopicArn),
	})
	return err
}
