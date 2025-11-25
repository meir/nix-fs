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
						Destination: "B",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						Destination: "B",
					},
				},
			},

			Expected: []Action{
				{
					Action: CREATE, // CREATE because file A obviously does not exist
					Location: Location{
						Origin:      "A",
						Destination: "B",
					},
				},
			},
		},
		{ // 4
			New: &StateFile{
				Locations: []Location{
					{
						Origin:      ".gitignore",
						Destination: ".gitignore",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{
					{
						Origin:      ".gitignore",
						Destination: ".gitignore",
					},
				},
			},

			Expected: []Action{
				{
					Action: CREATE, // CREATE because .gitignore exists, but is not a symlink
					Location: Location{
						Origin:      ".gitignore",
						Destination: ".gitignore",
					},
				},
			},
		},
		{ // 5
			New: &StateFile{
				Locations: []Location{
					{
						Origin:      "B",
						Destination: "B",
					},
				},
			},
			Old: &StateFile{
				Locations: []Location{
					{
						Origin:      "A",
						Destination: "B",
					},
				},
			},

			Expected: []Action{
				{
					Action: DELETE,
					Location: Location{
						Origin:      "A",
						Destination: "B",
					},
				},
				{
					Action: CREATE,
					Location: Location{
						Origin:      "B",
						Destination: "B",
					},
				},
			},
		},
	}

	for id, c := range cases {
		actions, err := Compare(c.New, c.Old)
		if err != nil {
			t.Errorf("[%d] Unexpected error: %s", id, err.Error())
		}

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
