package mytest

import (
	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/MeenaAlfons/go-shot/test-me/interfaces"
	"github.com/google/uuid"
)

type ExpectedNotification structs.Notification
type MessageGenerator interface {
	Name() string
	Start(appId string, count int) (<-chan ExpectedNotification, error)
}

type MessageGeneratorConfig struct {
	Queue            interfaces.Queue
	MaxBatchInterval int
	MaxBatchSize     int
}

func GenerateMessages(appId string, count int) []structs.Message {
	messages := make([]structs.Message, count)
	for i := 0; i < count; i++ {
		messages[i] = structs.Message{
			AppId: appId,
			Data:  uuid.NewString(),
		}
	}
	return messages
}
