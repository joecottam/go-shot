package mytest

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/google/uuid"
)

type MultipleAppsTestImpl struct {
	numApps          int
	name             string
	messageGenerator MessageGenerator
}

func NewMultipleAppsTest(numApps int, messageGenerator MessageGenerator) Test {
	return &MultipleAppsTestImpl{
		numApps:          numApps,
		name:             fmt.Sprintf("MultipleAppsTest(numApps=%d, messageGenerator=%s)", numApps, messageGenerator.Name()),
		messageGenerator: messageGenerator,
	}
}

func (t *MultipleAppsTestImpl) Name() string {
	return t.name
}

func (t *MultipleAppsTestImpl) Run(testConfig TestConfig) error {
	// Start message generators
	appIds := make([]string, t.numApps)
	notificationChannels := make(map[string]<-chan structs.Notification, t.numApps)
	expectedNotificationsChannels := make(map[string]<-chan ExpectedNotification, t.numApps)
	for i := 0; i < t.numApps; i++ {
		appIds[i] = uuid.NewString()
		notificationChannel, err := testConfig.SnsReceiver.Subscribe(appIds[i])
		if err != nil {
			return err
		}
		notificationChannels[appIds[i]] = notificationChannel
		expectedNotificationsChannels[appIds[i]], _ = t.messageGenerator.Start(appIds[i], 3)
	}

	// Wait for all messages to be received
	var wg sync.WaitGroup
	wg.Add(t.numApps)
	errorsChannel := make(chan error, t.numApps)
	for i := 0; i < t.numApps; i++ {
		go func(i int) {
			defer wg.Done()
			err := MatchNotifications(expectedNotificationsChannels[appIds[i]], notificationChannels[appIds[i]], time.Duration(testConfig.MaxBatchInterval+1)*time.Second)
			if err != nil {
				errorsChannel <- err
			}
		}(i)
	}
	wg.Wait()
	close(errorsChannel)

	errorsSlice := make([]error, 0)
	for err := range errorsChannel {
		errorsSlice = append(errorsSlice, err)
	}
	if len(errorsSlice) > 0 {
		errorMsg := fmt.Sprintf("%s failed with errors:", t.name)
		for _, err := range errorsSlice {
			errorMsg += "\n" + err.Error()
		}
		return errors.New(errorMsg)
	}

	return nil
}
