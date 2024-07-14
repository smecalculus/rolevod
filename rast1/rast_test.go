package rast

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// channels
	x channel = "x"
	y channel = "y"
	z channel = "z"

	// labels
	k label = "k"
)

var (
	// exps
	closeX = Close{x}

	// values
	closeV = CloseV{}
)

func TestEvalSucc(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		eta      eta
		exp      exp
		z        channel
		expected value
	}{
		"Id1":    {eta{y: closeV}, Id{x, y}, x, closeV},
		"Id2":    {eta{y: closeV}, Id{x, y}, z, closeV},
		"Lab1":   {eta{}, Lab{x, k, nil}, x, LabV{k, nil}},
		"Lab2":   {eta{x: CloCase{eta{}, branches{k: nil}, z}}, Lab{x, k, nil}, z, nil},
		"Case1":  {eta{}, Case{x, branches{k: nil}}, x, CloCase{eta{}, branches{k: nil}, x}},
		"Case2":  {eta{x: LabV{k, nil}}, Case{x, branches{k: nil}}, z, nil},
		"Send1":  {eta{y: nil}, Send{x, y, nil}, x, SendV{nil, nil}},
		"Send2":  {eta{x: CloRecv{eta{}, chan_exp{x, nil}, z}}, Send{x, y, nil}, z, nil},
		"Recv1":  {eta{}, Recv{x, y, nil}, x, CloRecv{eta{}, chan_exp{y, nil}, x}},
		"Recv2":  {eta{x: SendV{nil, nil}}, Recv{x, y, nil}, z, nil},
		"Close1": {eta{}, closeX, x, closeV},
		"Close2": {eta{}, closeX, z, closeV},
		"Wait1":  {eta{x: closeV}, Wait{x, nil}, x, nil},
		"Wait2":  {eta{x: closeV}, Wait{x, nil}, z, nil},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			actual, err := eval(env{}, tc.eta, tc.exp, tc.z)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(actual, tc.expected) {
				t.Errorf("\nactual: %T,\nexpected: %T,\ndiff: %s\n",
					actual, tc.expected, cmp.Diff(actual, tc.expected))
			}
		})
	}
}

func TestEvalErr(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		eta      eta
		exp      exp
		expected error
	}{
		"Wait": {eta{}, Wait{x, nil}, ErrUnexpectedWait},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			_, actual := eval(env{}, tc.eta, tc.exp, z)
			if !errors.Is(actual, tc.expected) {
				t.Errorf("\nactual: %#v,\nexpected: %#v",
					actual, tc.expected)
			}
		})
	}
}

func TestFoo(t *testing.T) {
	var foo exp = nil
	switch foo.(type) {
	case nil:
		t.Errorf("nil")
	}
}