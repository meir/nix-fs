package state

import (
	"testing"
)

func TestStateCompare(t *testing.T) {
	cases := []struct {
		New *StateFile
		Old *StateFile

		Expected []Action
	}{
		{ // 0
			New: &StateFile{},
			Old: &StateFile{},

			Expected: []Action{},
		},
		{ // 1
			New: &StateFile{
				Locations: []Location{},
			},
			Old: &StateFile{
				Locations: []Location{},
			},

			Expected: []Action{},
		},
		{ // 2
			New: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{},
			},

			Expected: []Action{
				{
					Action: CREATE,
					Location: Location{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
		},
		{ // 3
			New: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},

			Expected: []Action{
				{
					Action: NOOP,
					Location: Location{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
		},
		{ // 4
			New: &StateFile{
				Locations: []Location{
					{
						Origin:      "B",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},

			Expected: []Action{
				{
					Action: DELETE,
					Location: Location{
						Origin:      "A",
						IsDirectory: true,
						Destination: "B",
					},
				},
				{
					Action: CREATE,
					Location: Location{
						Origin:      "B",
						IsDirectory: true,
						Destination: "B",
					},
				},
			},
		},
	}

	for id, c := range cases {
		actions := Compare(c.New, c.Old)

		if len(actions) != len(c.Expected) {
			t.Errorf("[%d] len(actions) '%d' != len(expected) '%d'", id, len(actions), len(c.Expected))
		}

		for i, action := range actions {
			expect := c.Expected[i]

			if action.Action != expect.Action {
				t.Errorf("[%d] Given Action '%d' does not match expected Action '%d'", id, action.Action, expect.Action)
			}

			if !action.Location.Compare(expect.Location) {
				t.Errorf("[%d] Given Location '%s' -> '%s', does not match expected Location '%s' -> '%s'", id,
					action.Location.Origin,
					action.Location.Destination,

					expect.Location.Origin,
					expect.Location.Destination,
				)
			}
		}
	}
}
