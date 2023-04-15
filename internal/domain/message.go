package domain

const (
	DirectMessageType = "directMessage"
	GroupMessageType  = "groupMessage"
)

type Message struct {
	Type     string
	Sender   string
	Group    string
	Receiver string
	Message  string
}
