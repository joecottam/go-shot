package mytest

import (
	"errors"
	"reflect"
	"time"

	"github.com/MeenaAlfons/go-shot/structs"
)

func MatchNotifications(expectedNotificationsChannel <-chan ExpectedNotification, actualNotificationsChannel <-chan structs.Notification, notificationTimeout time.Duration) error {
	for notification := range expectedNotificationsChannel {
		select {
		case <-time.After(notificationTimeout):
			return errors.New("Timed out waiting for notification")
		case receivedNotification := <-actualNotificationsChannel:
			if len(receivedNotification.Messages) != len(notification.Messages) {
				return errors.New("Incorrect number of messages received")
			}
			if !reflect.DeepEqual(receivedNotification.Messages, notification.Messages) {
				return errors.New("Incorrect messages received")
			}
		}
	}
	return nil
}
