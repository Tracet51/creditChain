package voter

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/google/uuid"
)

const (
	voterRecordFilename = "/tmp/voter_record.txt"
)

// Voter is responsible to voting
type Voter struct {
	ID        uuid.UUID
	QuorumSet map[string]*Voter
}

// GetVoter creates a new node
func GetVoter() *Voter {

	var voter Voter
	data, err := ioutil.ReadFile(voterRecordFilename)

	if err != nil {
		voter = createVoterRecord()
	} else {
		err = json.Unmarshal(data, &voter)
	}
	
	return &voter
}

func createVoterRecord() Voter {

	id, _ := uuid.NewUUID()
	voter := Voter{
		ID:        id,
		QuorumSet: make(map[string]*Voter)}

	voterData, err := json.Marshal(voter)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(voterRecordFilename, voterData, 0667)
	if err != nil {
		log.Fatalln("Could not create record of node: ", err)
	}

	return voter
}