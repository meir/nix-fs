package state

import (
	"os"

	"golang.org/x/exp/slices"
)

type LinkAction uint8

const (
	DELETE LinkAction = iota
	NOOP
	CREATE
)

type Action struct {
	Action   LinkAction
	Location Location
}

func sortActions(a Action, b Action) int {
	return int(a.Action) - int(b.Action)
}

func Compare(new *StateFile, old *StateFile) ([]Action, error) {
	var actions []Action

CreateOrNOOP:
	for _, locationNew := range new.Locations {
		for i, locationOld := range old.Locations {
			if locationNew.Compare(locationOld) {
				action := NOOP

				stat, err := os.Lstat(locationNew.Destination)
				if err != nil && !os.IsNotExist(err) {
					return nil, err
				} else if os.IsNotExist(err) || stat.Mode()&os.ModeSymlink == 0 {
					action = CREATE
				}

				actions = append(actions, Action{
					Action:   action,
					Location: locationNew,
				})

				// remove matched old location to optimize further searches and delete remaining paths
				old.Locations = append(old.Locations[:i], old.Locations[i+1:]...)
				continue CreateOrNOOP
			}
		}

		actions = append(actions, Action{
			Action:   CREATE,
			Location: locationNew,
		})
	}

	for _, locationOld := range old.Locations {
		actions = append(actions, Action{
			Action:   DELETE,
			Location: locationOld,
		})
	}

	slices.SortFunc(actions, sortActions)
	return actions, nil
}
