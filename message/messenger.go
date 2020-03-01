package message

type Messenger struct {
	InboundMessages  chan Message
	OutboundMessages chan Message
}
