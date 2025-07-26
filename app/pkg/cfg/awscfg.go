package cfg

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type AWSService struct {
	SNSClient *sns.Client
	Config    *Config
}

func NewAWSService(cfg *Config) (*AWSService, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.AWS.Region))
	if err != nil {
		return nil, err
	}

	var snsClient *sns.Client
	if cfg.IsLocalStack() {
		snsClient = sns.NewFromConfig(awsCfg, func(o *sns.Options) {
			o.BaseEndpoint = aws.String(cfg.AWS.LocalStack.Endpoint)
		})
	} else {
		snsClient = sns.NewFromConfig(awsCfg)
	}

	return &AWSService{
		SNSClient: snsClient,
		Config:    cfg,
	}, nil
}

func (a *AWSService) GetSNSClient() *sns.Client {
	return a.SNSClient
}

func (a *AWSService) GetSNSTopicArn() string {
	return a.Config.GetSNSTopicArn()
}
