package message

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
)

type MessageCodes map[int]reflect.Type
type MessageChannels map[reflect.Type]chan Message

var messageCodes MessageCodes
var messageChannels MessageChannels

func Register(messageType interface{}, messageChannel chan Message) {

	if messageCodes == nil {
		messageCodes = make(MessageCodes, 0)
	}

	if messageChannels == nil {
		messageChannels = make(MessageChannels, 0)
	}

	nextMessageCode := getNextMessageCode()

	typeElement := reflect.TypeOf((*InitMessage)(nil)).Elem()
	messageCodes[nextMessageCode] = typeElement
	messageChannels[typeElement] = messageChannel
}

func getNextMessageCode() int {
	maxKey := 0
	for key := range messageCodes {
		if key > maxKey {
			maxKey = key
		}
	}

	return maxKey
}

func GetMessage(message []byte) {

	for key := range messageCodes {
		if bytes.IndexByte(message, byte(key)) == 0 {
			log.Printf(string(message))
			messageType := messageCodes[key]
			payload := reflect.New(messageType)
			err := json.Unmarshal(message, &payload)
			if err != nil {
				log.Fatal(err.Error())
			}

			messageChannels[messageType] <- payload
		}
	}
}

func PrepareMessageForSending(messageCode int, message []byte) []byte {

	preparedMessage := prependMessageCode(messageCode, message)
	return preparedMessage
}

func prependMessageCode(messageCode int, data []byte) []byte {
	codeSlice := []byte(string(messageCode))
	prependedBytes := append(codeSlice, data...)
	return prependedBytes
}
