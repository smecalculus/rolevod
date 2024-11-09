package role_test

import (
	"os"
	"testing"

	"smecalculus/rolevod/app/role"
)

var (
	api = role.NewAPI()
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
