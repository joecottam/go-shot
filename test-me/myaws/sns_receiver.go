package myaws

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/MeenaAlfons/go-shot/structs"
	"github.com/MeenaAlfons/go-shot/test-me/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsRecieverConfig struct {
	Port      int
	Host      string
	AwsConfig aws.Config
}

type SnsRecieverImpl struct {
	config   SnsRecieverConfig
	server   *http.Server
	client   *sns.Client
	channels map[string]chan structs.Notification

	subscriptionConfirmationTokenChannels map[string]chan string
}

func NewSnsReceiver(config SnsRecieverConfig) (interfaces.SnsReceiver, error) {
	s := &SnsRecieverImpl{
		config:                                config,
		client:                                sns.NewFromConfig(config.AwsConfig),
		channels:                              make(map[string]chan structs.Notification),
		subscriptionConfirmationTokenChannels: make(map[string]chan string),
	}
	server, err := NewServer(s.config.Port, http.HandlerFunc(s.handler))
	if err != nil {
		return nil, err
	}
	s.server = server
	return s, nil
}

func (s *SnsRecieverImpl) Subscribe(appId string) (chan structs.Notification, error) {
	if channel, ok := s.channels[appId]; ok {
		return channel, nil
	}

	// Create a topic
	output, err := s.client.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name: aws.String(appId),
	})
	if err != nil {
		return nil, err
	}

	// Subscribe to the topic
	s.subscriptionConfirmationTokenChannels[appId] = make(chan string)
	_, err = s.client.Subscribe(context.TODO(), &sns.SubscribeInput{
		Endpoint: aws.String(fmt.Sprintf("http://%s:%d/%s", s.config.Host, s.config.Port, appId)),
		Protocol: aws.String("http"),
		TopicArn: output.TopicArn,
	})
	if err != nil {
		return nil, err
	}

	// Confirm the subscription
	token := <-s.subscriptionConfirmationTokenChannels[appId]
	_, err = s.client.ConfirmSubscription(context.TODO(), &sns.ConfirmSubscriptionInput{
		Token:    aws.String(token),
		TopicArn: output.TopicArn,
	})
	if err != nil {
		return nil, err
	}

	// Create a channel
	s.channels[appId] = make(chan structs.Notification, 10)
	return s.channels[appId], nil
}

func NewServer(port int, handler http.Handler) (*http.Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		// Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	go func() {
		server.Serve(l)
	}()
	return server, nil
}

func (s *SnsRecieverImpl) handler(w http.ResponseWriter, r *http.Request) {
	appId := r.URL.Path[1:]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var msg struct {
		Type    string `json:"Type"`
		Token   string `json:"Token"`
		Message string `json:"Message"`
	}
	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		log.Println(err)
	} else {
		switch msg.Type {
		case "SubscriptionConfirmation":
			if _, ok := s.subscriptionConfirmationTokenChannels[appId]; !ok {
				log.Println("Subscription confirmation for unknown app", appId)
				break
			}
			s.subscriptionConfirmationTokenChannels[appId] <- msg.Token
		case "Notification":
			if _, ok := s.channels[appId]; !ok {
				log.Println("Notification for unknown app", appId)
				break
			}
			var notification structs.Notification
			err := json.Unmarshal([]byte(msg.Message), &notification)
			switch err {
			case nil:
				s.channels[appId] <- notification
			default:
				log.Println("Can't parse notification", err)
			}
		}
	}

	w.Write([]byte("OK"))
}
