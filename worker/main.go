package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MeenaAlfons/go-shot/config"
	"github.com/MeenaAlfons/go-shot/localstack"
	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	c := make(chan structs.Message)
	go getMessages(c, awsConfig)

	for {
		msg := <-c
		fmt.Println("Message: ", msg)
	}
}

func getMessages(c chan structs.Message, awsConfig *aws.Config) {

	cfg, err := config.GetConfig()
	sqsService := sqs.NewFromConfig(*awsConfig)
	ctx := context.TODO()

	if err != nil {
		fmt.Println("Error from cfg: ", err)
	}

	for {
		msgOutput, err := sqsService.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl: &cfg.Queue,
		})
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			for _, msg := range msgOutput.Messages {
				message := structs.Message{}
				err = json.Unmarshal([]byte(*msg.Body), &message)
				if err != nil {
					fmt.Println("Error unmarshalling json: ", err)
				}
				c <- message
			}
		}
	}
}
