package mytest

import (
	"log"
	"math/rand"
	"time"
)

type ScenarioMessageGeneratorImpl struct {
	config          MessageGeneratorConfig
	scenarioFactory ScenarioFactory
}

type ScenarioFactory interface {
	Name() string
	GenerateScenario(maxBatchInterval, maxBatchSize int) (int, time.Duration)
}

func NewScenarioMessageGenerator(config MessageGeneratorConfig, scenarioFactory ScenarioFactory) MessageGenerator {
	return &ScenarioMessageGeneratorImpl{
		config:          config,
		scenarioFactory: scenarioFactory,
	}
}

func (m *ScenarioMessageGeneratorImpl) Name() string {
	return m.scenarioFactory.Name()
}

func (m *ScenarioMessageGeneratorImpl) Start(appId string, count int) (<-chan ExpectedNotification, error) {
	channel := make(chan ExpectedNotification, count)
	go func() {
		for i := 0; i < count; i++ {
			count, durationBetweenMessages := m.scenarioFactory.GenerateScenario(m.config.MaxBatchInterval, m.config.MaxBatchSize)
			log.Printf("Sending %d messages with %s interval, appId %s", count, durationBetweenMessages, appId)

			// Generate messages
			messages := GenerateMessages(appId, count)
			channel <- ExpectedNotification{
				Messages: messages,
			}

			// Send messages
			for _, message := range messages {
				err := m.config.Queue.SendMessage(message)
				if err != nil {
					panic(err)
				}
				<-time.After(durationBetweenMessages)
			}

			// For batch interval scenario, we want to make sure that
			// the batch interval is passed before sending the next batch.
			// At this point, we are exacly at the end of the batch interval.
			// Wait for a bit to cross the batch interval boundary.
			<-time.After(200 * time.Millisecond)
		}
		close(channel)
	}()
	return channel, nil
}

// BatchIntervalScenario:
// MaxBatchInterval, num of messages, and the interval between messages
// are related by the following equation:
// interval between messages = MaxBatchInterval / num of messages
// where num of messages is between 1 and MaxBatchSize-1
type BatchIntervalScenarioFactory struct{}

func NewBatchIntervalScenarioFactory() ScenarioFactory {
	return &BatchIntervalScenarioFactory{}
}
func (f *BatchIntervalScenarioFactory) Name() string {
	return "BatchIntervalScenario"
}
func (f *BatchIntervalScenarioFactory) GenerateScenario(maxBatchInterval, maxBatchSize int) (int, time.Duration) {
	// Generate a random number of messages to send
	// between 1 and maxBatchSize-1
	count := 1 + rand.Intn(maxBatchSize-1)

	// Calculate the interval between messages
	// to make sure that the total time to send all messages
	// is less than maxBatchInterval
	duration := time.Millisecond * time.Duration(1000*maxBatchInterval/count)
	return count, duration
}

// BatchSizeScenario:
// MaxBatchSize, num of messages, and the interval between messages
// are related by the following equation:
// num of messages = MaxBatchSize
// interval between messages < MaxBatchInterval / MaxBatchSize
type BatchSizeScenarioFactory struct{}

func NewBatchSizeScenarioFactory() ScenarioFactory {
	return &BatchSizeScenarioFactory{}
}
func (f *BatchSizeScenarioFactory) Name() string {
	return "BatchSizeScenario"
}
func (f *BatchSizeScenarioFactory) GenerateScenario(maxBatchInterval, maxBatchSize int) (int, time.Duration) {
	maxDuration := time.Millisecond * time.Duration(1000*maxBatchInterval/maxBatchSize)
	minDuration := maxDuration / 2
	duration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration)))
	return maxBatchSize, duration
}
