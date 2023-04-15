package main

import (
	"log"

	"github.com/MeenaAlfons/go-shot/config"
	"github.com/MeenaAlfons/go-shot/localstack"
	"github.com/MeenaAlfons/go-shot/test-me/myaws"
	"github.com/MeenaAlfons/go-shot/test-me/mytest"
)

func main() {
	log.Println("Starting test runner...")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	awsConfig, err := localstack.AwsConfigFromEndpoint(cfg.AwsEndpoint, cfg.AwsRegion)
	if err != nil {
		log.Fatal(err)
	}

	snsReceiver, err := myaws.NewSnsReceiver(myaws.SnsRecieverConfig{
		AwsConfig: *awsConfig,
		Host:      cfg.Host,
		Port:      cfg.Port,
	})
	if err != nil {
		log.Fatal(err)
	}

	queue, err := myaws.NewQueue(*awsConfig, cfg.Queue)
	if err != nil {
		log.Fatal(err)
	}

	// Start the test runner
	// which will put messages on SQS
	// and then expect to receive batched messages on SNS
	testConfig := mytest.TestConfig{
		SnsReceiver:      snsReceiver,
		Queue:            queue,
		MaxBatchInterval: cfg.MaxBatchInterval,
		MaxBatchSize:     cfg.MaxBatchSize,
	}

	messageGeneratorConfig := mytest.MessageGeneratorConfig{
		Queue:            queue,
		MaxBatchInterval: cfg.MaxBatchInterval,
		MaxBatchSize:     cfg.MaxBatchSize,
	}
	batchIntervalScenario := mytest.NewScenarioMessageGenerator(messageGeneratorConfig, mytest.NewBatchIntervalScenarioFactory())
	batchSizeScenario := mytest.NewScenarioMessageGenerator(messageGeneratorConfig, mytest.NewBatchSizeScenarioFactory())

	runner := mytest.NewTestRunner(testConfig, []mytest.Test{
		mytest.NewSynchronizeTest(),
		mytest.NewMultipleAppsTest(1, batchIntervalScenario),
		mytest.NewMultipleAppsTest(3, batchIntervalScenario),
		mytest.NewMultipleAppsTest(1, batchSizeScenario),
		mytest.NewMultipleAppsTest(3, batchSizeScenario),
	})
	err = runner.Run()

	switch err {
	case nil:
		log.Println("All tests passed")
	default:
		log.Fatal(err)
	}
}
