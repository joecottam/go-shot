package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MeenaAlfons/go-shot/config"
	"github.com/MeenaAlfons/go-shot/localstack"
	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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

	messages := make(map[string][]structs.Message)

	snsService := sns.NewFromConfig(*awsConfig)

	for {
		msg := <-c
		messages[msg.AppId] = append(messages[msg.AppId], msg)

		if len(messages[msg.AppId]) == cfg.MaxBatchSize {
			fmt.Println("Batch size reached for app: ", msg.AppId)
			publishAndClear(messages, msg, snsService)
		}

		if len(messages[msg.AppId]) == 1 {
			go func() {
				// As soon as we have the first messages we sleep for the time set
				<-time.After(time.Duration(cfg.MaxBatchInterval) * time.Second)
				if len(messages[msg.AppId]) > 0 {
					fmt.Println("Max batch interval reached for app: ", msg.AppId)
					publishAndClear(messages, msg, snsService)
				}
			}()
		}
	}
}

func publishAndClear(messages map[string][]structs.Message, msg structs.Message, snsService *sns.Client) {
	notification := structs.Notification{
		Messages: messages[msg.AppId],
	}

	notificationJson, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
	}

	topic, err := snsService.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name: aws.String(msg.AppId),
	})

	if err != nil {
		fmt.Println("Error creating topic: ", err)
	}

	_, err = snsService.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(string(notificationJson)),
		TopicArn: topic.TopicArn,
	})
	if err != nil {
		fmt.Println("Error publishing message: ", err)
	}
	messages[msg.AppId] = []structs.Message{}
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
				if err = json.Unmarshal([]byte(*msg.Body), &message); err != nil {
					fmt.Println("Error unmarshalling json: ", err)
				}
				c <- message
				_, err = sqsService.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &cfg.Queue,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					fmt.Println("Error deleting message: ", err)
				}
			}
		}
	}
}
