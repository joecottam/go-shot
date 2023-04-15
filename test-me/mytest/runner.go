package mytest

import (
	"errors"
	"log"

	"github.com/MeenaAlfons/go-shot/test-me/interfaces"
)

type Test interface {
	Run(TestConfig) error
	Name() string
}

type TestRunner interface {
	Run() error
}

type TestConfig struct {
	Queue            interfaces.Queue
	SnsReceiver      interfaces.SnsReceiver
	MaxBatchInterval int
	MaxBatchSize     int
}

func NewTestRunner(testConfig TestConfig, tests []Test) TestRunner {
	return &TestRunnerImpl{
		testConfig: testConfig,
		tests:      tests,
	}
}

type TestRunnerImpl struct {
	testConfig TestConfig
	tests      []Test
}

func (t *TestRunnerImpl) Run() error {
	var occuredErrors []error
	for _, test := range t.tests {
		log.Println("Running test:", test.Name())
		if err := test.Run(t.testConfig); err != nil {
			log.Printf("%s failed with error: %v", test.Name(), err)
			occuredErrors = append(occuredErrors, err)
		}
	}
	if len(occuredErrors) > 0 {
		return errors.New("There were errors running tests")
	}
	return nil
}
