package sig_test

import (
	"os"
	"testing"

	"smecalculus/rolevod/app/sig"
)

var (
	api = sig.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
