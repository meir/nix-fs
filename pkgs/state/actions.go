package state

import "golang.org/x/exp/slices"

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

func Compare(new *StateFile, old *StateFile) []Action {
	var actions []Action

CreateOrNOOP:
	for _, locationNew := range new.Locations {
		for i, locationOld := range old.Locations {
			if locationNew.Compare(locationOld) {
				actions = append(actions, Action{
					Action:   NOOP,
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
	return actions
}
