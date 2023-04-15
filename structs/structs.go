package structs

type Message struct {
	AppId string `json:"AppId"`
	Data  string `json:"Data"`
}

type Notification struct {
	Messages []Message `json:"Messages"`
}
