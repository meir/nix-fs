package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const STATE_VERSION = 1

type StateFile struct {
	StateVersion   uint8      `json:"version"`
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
	actions, err := Compare(sf, old)
	if err != nil {
		return err
	}

	for _, action := range actions {
		switch action.Action {
		case DELETE:
			fmt.Printf("Deleting link: %s -> %s\n", action.Location.Origin, action.Location.Destination)
			action.Location.RemoveLink()
		case CREATE:
			fmt.Printf("Creating link: %s -> %s\n", action.Location.Origin, action.Location.Destination)
			action.Location.CreateLink()
		case NOOP:
			fmt.Printf("No action for link: %s -> %s\n", action.Location.Origin, action.Location.Destination)
		}
	}

	return nil
}
