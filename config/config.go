package config

import (
	env "github.com/caarlos0/env/v8"
)

type Config struct {
	AwsEndpoint string `env:"AWS_ENDPOINT,required"`
	AwsRegion   string `env:"AWS_REGION,required"`

	Queue            string `env:"QUEUE,required"`
	Topic            string `env:"TOPIC,required"`
	MaxBatchInterval int    `env:"MAX_BATCH_INTERVAL,required"`
	MaxBatchSize     int    `env:"MAX_BATCH_SIZE,required"`

	// These are for the test runner
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
