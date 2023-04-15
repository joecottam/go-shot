package myaws

import (
	"context"
	"encoding/json"

	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/MeenaAlfons/go-shot/test-me/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueImpl struct {
	client   *sqs.Client
	QueueUrl string
}

func NewQueue(awsConfig aws.Config, name string) (interfaces.Queue, error) {
	// Create a SQS service client.
	client := sqs.NewFromConfig(awsConfig)

	// Create a queue.
	result, err := client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	return &QueueImpl{
		client:   client,
		QueueUrl: *result.QueueUrl,
	}, nil
}

func (q *QueueImpl) SendMessage(message structs.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = q.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(data)),
		QueueUrl:    &q.QueueUrl,
	})
	return err
}

// func CreateQueue(awsConfig aws.Config, name string) (string, error) {
// 	// Create a SQS service client.
// 	svc := sqs.NewFromConfig(awsConfig)

// 	// Create a queue.
// 	result, err := svc.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
// 		QueueName: aws.String(name),
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return *result.QueueUrl, nil
// }
