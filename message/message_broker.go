package message

type MessageBroker struct {
	InboundMessage  chan Message
	OutboundMessage chan Message
}
