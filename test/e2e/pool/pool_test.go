package pool_test

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/app/pool"
)

var (
	api = pool.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreation(t *testing.T) {

	t.Run("CreateRetreive", func(t *testing.T) {
		// given
		poolSpec1 := pool.Spec{Title: "ts1"}
		poolRoot1, err := api.Create(poolSpec1)
		if err != nil {
			t.Fatal(err)
		}
		// and
		poolSpec2 := pool.Spec{Title: "ts2", SupID: poolRoot1.ID}
		poolRoot2, err := api.Create(poolSpec2)
		if err != nil {
			t.Fatal(err)
		}
		// when
		poolSnap1, err := api.Retrieve(poolRoot1.ID)
		if err != nil {
			t.Fatal(err)
		}
		// then
		extectedSub := pool.ConvertRootToRef(poolRoot2)
		if !slices.Contains(poolSnap1.Subs, extectedSub) {
			t.Errorf("unexpected subs in %q; want: %+v, got: %+v",
				poolSpec1.Title, extectedSub, poolSnap1.Subs)
		}
	})
}
