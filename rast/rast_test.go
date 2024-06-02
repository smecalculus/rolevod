package rast

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const x channel = "x"
const z channel = "z"

func TestEvalSuccess(t *testing.T) {
	t.Parallel()
	var closeX = Close{x}
	var closeV = CloseV{}
	tcs := map[string]struct {
		eta      eta
		exp      exp
		expected value
	}{
		"Close": {eta{}, closeX, closeV},
		"Wait":  {eta{x: closeV}, Wait{x, nil}, nil},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			actual, err := eval(env{}, tc.eta, tc.exp, z)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(actual, tc.expected) {
				t.Errorf("\nactual: %T,\nexpected: %T,\ndiff: %s\n", actual, tc.expected, cmp.Diff(actual, tc.expected))
			}
		})
	}
}

func TestEvalError(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		eta      eta
		exp      exp
		expected error
	}{
		"Wait":  {eta{}, Wait{x, nil}, ErrUnexpectedWait},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			_, actual := eval(env{}, tc.eta, tc.exp, z)
			if !errors.Is(actual, tc.expected) {
				t.Errorf("\nactual: %#v,\nexpected: %#v", actual, tc.expected)
			}
		})
	}
}
