package step

import (
	"os"
	"slices"
	"testing"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCollectCtx(t *testing.T) {
	// given
	chID := id.New()
	// and
	term := SpawnSpec{Ctx: []chnl.ID{chID}, Cont: CloseSpec{}}
	// when
	actualCtx := CollectCtx(id.New(), term)
	// then
	if !slices.Contains(actualCtx, chID) {
		t.Errorf("unexpected ctx: want %q in %v", chID, actualCtx)
	}
}
