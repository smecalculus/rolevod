package team_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/app/team"
)

var (
	api = team.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreation(t *testing.T) {

	t.Run("CreateRetreive", func(t *testing.T) {
		// given
		teamSpec1 := team.Spec{Title: "ts1"}
		teamRoot1, err := api.Create(teamSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		teamSpec2 := team.Spec{Title: "ts2", SupID: teamRoot1.ID}
		teamRoot2, err := api.Create(teamSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// when
		teamSnap1, err := api.Retrieve(teamRoot1.ID)
		if err != nil {
			t.Fatal(err)
		}
		// then
		extectedSub := team.ConvertRootToRef(teamRoot2)
		if !slices.Contains(teamSnap1.Subs, extectedSub) {
			t.Errorf("unexpected subs in %q; want: %+v, got: %+v",
				teamSpec1.Title, extectedSub, teamSnap1.Subs)
		}
	})
}
