package message

import (
	"encoding/json"
	"log"
	"time"
)

type Message interface {
}

type InitMessage struct {
	To     string
	From   string
	Amount float64
	Date   time.Time
	Hash   string
}

type MessageInt struct {
	Payload int
}

type MessageString struct {
	Payload string
}

// ToBytes Converts a Message to Bytes
func ToBytes(message *Message) []byte {
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}
