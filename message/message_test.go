package message

import (
	"encoding/json"
	"testing"

	"github.com/tracet51/creditChain/message"
)

func TestToBytes(t *testing.T) {
	initMessage := &message.InitMessage{}
	data := message.ToBytes(initMessage)
	if data == nil {
		t.Error("Ahh")
	}
}

func TestFromBytes(t *testing.T) {
	initMessage := message.InitMessage{
		To:     "Trace",
		From:   "Kylie",
		Amount: 50,
	}
	data, _ := json.Marshal(initMessage)
	newMessage := message.FromBytes(data)

	if newMessage == nil {
		t.Error("AHHH")
	}

	initMessage, ok := newMessage.(message.InitMessage)

	if ok == false {
		t.Error(ok)
	}

}
