package mytest

import (
	"errors"
	"log"
	"time"

	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/google/uuid"
)

type SynchronizeTestImpl struct {
}

func NewSynchronizeTest() Test {
	return &SynchronizeTestImpl{}
}

func (t *SynchronizeTestImpl) Name() string {
	return "SynchronizeTest"
}

func (t *SynchronizeTestImpl) Run(testConfig TestConfig) error {
	err := t.synchronize(testConfig, structs.Message{AppId: "synch-1", Data: uuid.NewString()})
	if err != nil {
		return err
	}

	// Synchronize one more time
	err = t.synchronize(testConfig, structs.Message{AppId: "synch-2", Data: uuid.NewString()})
	return err
}

func (t *SynchronizeTestImpl) synchronize(testConfig TestConfig, message structs.Message) error {
	notificationChannel, err := testConfig.SnsReceiver.Subscribe(message.AppId)
	if err != nil {
		return err
	}

	// Put synchronization message on queue
	err = testConfig.Queue.SendMessage(message)
	if err != nil {
		return err
	}

	// Wait for message, skipping any other messages
	for {
		select {
		case received := <-notificationChannel:
			if len(received.Messages) == 1 && received.Messages[0] == message {
				log.Println("Received synchronization message", message)
				return nil
			}
			log.Println("Skipping message:", received)

		// Wait for more than the maximum batch interval to ensure we receive this message alone in a batch
		case <-time.After(time.Duration(1+testConfig.MaxBatchInterval) * time.Second):
			return errors.New("Timeout waiting for synchronization message")
		}
	}
}
