package interfaces

import "github.com/MeenaAlfons/go-shot/structs"

type SnsReceiver interface {
	Subscribe(appId string) (chan structs.Notification, error)
}
