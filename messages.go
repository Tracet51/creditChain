package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Payload string
	NodeID  string
}

type MessageInt struct {
	Payload int
}

type MessageString struct {
	Payload string
}

// MessageToBytes Converts a Message to Bytes
func MessageToBytes(message *Message) []byte {
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}
