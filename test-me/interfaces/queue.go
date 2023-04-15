package interfaces

import "github.com/MeenaAlfons/go-shot/structs"

type Queue interface {
	SendMessage(message structs.Message) error
}
