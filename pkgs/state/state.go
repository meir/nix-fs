package state

import (
	"fmt"
	"time"
)

type StateFile struct {
	LastGeneration time.Time  `json:"time"`
	Locations      []Location `json:"locations"`
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
