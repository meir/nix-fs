package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const STATE_VERSION = 1

type StateFile struct {
	StateVersion   uint8      `json:"state_version"`
	LastGeneration time.Time  `json:"time"`
	Locations      []Location `json:"locations"`
}

func NewStateFile(path string) (*StateFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var stateFile StateFile
	err = json.Unmarshal(data, &stateFile)
	if err != nil {
		return nil, err
	}

	return &stateFile, nil
}

func EmptyStateFile() *StateFile {
	return &StateFile{
		StateVersion:   STATE_VERSION,
		LastGeneration: time.Now(),
		Locations:      []Location{},
	}
}

func (sf *StateFile) Apply(old *StateFile) error {
	actions := Compare(sf, old)

	for _, action := range actions {
		switch action.Action {
		case DELETE:
			fmt.Println("Deleting link:", action.Location.Destination)
			action.Location.RemoveLink()
		case CREATE:
			fmt.Println("Creating link:", action.Location.Destination)
			action.Location.CreateLink()
		case NOOP:
			fmt.Println("No action for link:", action.Location.Destination)
		}
	}

	return nil
}
