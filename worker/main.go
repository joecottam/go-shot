package main

import (
	"log"

	"github.com/MeenaAlfons/go-shot/config"
	"github.com/MeenaAlfons/go-shot/localstack"
)

func main() {
	log.Println("Starting worker...")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	awsConfig, err := localstack.AwsConfigFromEndpoint(cfg.AwsEndpoint, cfg.AwsRegion)
	if err != nil {
		log.Fatal(err)
	}

	// Reding config and setting up awsConfig
	// is already done for you. You can use
	// the awsConfig variable to acess AWS services
	// like SQS, SNS, etc.
	// The following is a dummy line:
	_ = awsConfig
}
