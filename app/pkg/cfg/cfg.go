package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	AWS     AWSConfig     `mapstructure:"aws"`
	SNS     SNSConfig     `mapstructure:"sns"`
	Server  ServerConfig  `mapstructure:"server"`
	Logging LoggingConfig `mapstructure:"logging"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type AWSConfig struct {
	Region     string           `mapstructure:"region"`
	LocalStack LocalStackConfig `mapstructure:"localstack"`
}

type LocalStackConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

type SNSConfig struct {
	TopicArn     string `mapstructure:"topic_arn"`
	TopicArnProd string `mapstructure:"topic_arn_prod"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	config.overrideWithEnvVars()

	return &config, nil
}

func (c *Config) overrideWithEnvVars() {
	if env := os.Getenv("APP_ENVIRONMENT"); env != "" {
		c.App.Environment = env
	}

	if region := os.Getenv("AWS_REGION"); region != "" {
		c.AWS.Region = region
	}

	if enabled := os.Getenv("AWS_LOCALSTACK_ENABLED"); enabled != "" {
		c.AWS.LocalStack.Enabled = enabled == "true"
	}

	if endpoint := os.Getenv("AWS_LOCALSTACK_ENDPOINT"); endpoint != "" {
		c.AWS.LocalStack.Endpoint = endpoint
	}

	if topicArn := os.Getenv("SNS_TOPIC_ARN"); topicArn != "" {
		c.SNS.TopicArn = topicArn
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p := viper.GetInt("SERVER_PORT"); p > 0 {
			c.Server.Port = p
		}
	}

	if level := os.Getenv("LOGGING_LEVEL"); level != "" {
		c.Logging.Level = level
	}
}

func (c *Config) IsLocalStack() bool {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		return false
	}

	if useLocalStack := os.Getenv("USE_LOCALSTACK"); useLocalStack != "" {
		return useLocalStack == "true"
	}

	return c.AWS.LocalStack.Enabled
}

func (c *Config) GetSNSTopicArn() string {
	if c.IsLocalStack() || c.App.Environment == "local" {
		return c.SNS.TopicArn
	}
	return c.SNS.TopicArnProd
}
