package e2e

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/app/client"
	"smecalculus/rolevod/app/dcl"
	ws "smecalculus/rolevod/app/ws"
)

var (
	envApi = client.NewEnvApi()
	tpApi  = client.NewTpApi()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestTpIntro(t *testing.T) {
	// given
	es := ws.EnvSpec{Name: "env-1"}
	env, err := envApi.Create(es)
	if err != nil {
		t.Fatal(err)
	}
	// and
	ts := dcl.TpSpec{Name: "tp-1", St: dcl.One{}}
	tp, err := tpApi.Create(ts)
	if err != nil {
		t.Fatal(err)
	}
	// when
	intro := ws.TpIntro{EnvID: env.ID, TpID: tp.ID}
	err = envApi.Introduce(intro)
	if err != nil {
		t.Fatal(err)
	}
	// and
	actual, err := envApi.Retrieve(env.ID)
	if err != nil {
		t.Fatal(err)
	}
	// then
	expectedTp := dcl.TpTeaserFromTpRoot(tp)
	if !slices.Contains(actual.Tps, expectedTp) {
		t.Errorf("unexpected tps in %q; want: %+v, got: %+v", env.Name, expectedTp, actual.Tps)
	}
}
