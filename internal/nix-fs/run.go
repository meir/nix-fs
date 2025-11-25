package nixfs

import (
	"encoding/json"
	"os"
	"time"

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

	err = newState.Apply(oldState)
	if err != nil {
		return err
	}

	return OverwriteState(newState, oldStateFileLocation)
}

// OverwriteState overwrites the old state file with new one
func OverwriteState(newState *state.StateFile, stateLocation string) error {
	newState.LastGeneration = time.Now()
	data, err := json.MarshalIndent(newState, "", "  ")
	if err != nil {
		return nil
	}

	return os.WriteFile(stateLocation, data, 0644)
}
