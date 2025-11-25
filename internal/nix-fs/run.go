package nixfs

import (
	"encoding/json"
	"os"

	"github.com/meir/nix-fs/pkgs/state"
)

func getPreviousState(location string) (*state.StateFile, error) {
	data, err := os.ReadFile(location)
	if os.IsNotExist(err) {
		return state.EmptyStateFile(), nil
	} else if err != nil {
		return nil, err
	}

	var sf state.StateFile
	err = json.Unmarshal(data, &sf)
	if err != nil {
		return nil, err
	}

	return &sf, nil
}

func Run(stateFileLocation, oldStateFileLocation string) error {
	newState, err := state.NewStateFile(stateFileLocation)
	if err != nil {
		return err
	}

	oldState, err := getPreviousState(oldStateFileLocation)
	if err != nil {
		return err
	}

	return newState.Apply(oldState)
}
