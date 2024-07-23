package exec

import (
	"errors"
	a "smecalculus/rolevod/rast2/ast"
	"testing"

	// "github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestExecSucc(t *testing.T) {
	tc := func(env a.Environment, f a.Expname, expected *Configuration) {
		t.Helper()
		actual, err := Exec(env, f)
		if err != nil {
			t.Fatalf("unexpected error %q", err)
		}
		// deep.CompareUnexportedFields = true
		// deep.NilMapsAreEmpty = true
		// deep.NilSlicesAreEmpty = true
		// deep.NilPointersAreZero = true
		// if diff := deep.Equal(actual, expected); diff != nil {
		// 	t.Errorf("unexpected config:\n%v\n", strings.Join(diff, "\n"))
		// }
		opt := cmpopts.EquateEmpty()
		if diff := cmp.Diff(expected, actual, opt); diff != "" {
			t.Errorf("unexpected config (-want +got):\n%s", diff)
		}
	}

	foo := "foo"
	z := a.Chan{Id: "z"}
	main := "main"
	d := a.Chan{Id: "d"}
	env := map[string]a.Decl{
		foo: a.ExpDecDef{
			F: foo,
			Zc: a.ChanTp{
				X: z,
				A: a.One{},
			},
			P: a.Close{X: z},
		},
		main: a.ExpDecDef{
			F: main,
			Zc: a.ChanTp{
				X: d,
				A: a.One{},
			},
			P: a.Spawn{
				X: z,
				F: foo,
				Q: a.Wait{X: z, P: a.Close{X: d}},
			},
		},
	}
	config := &Configuration{
		Conf: map[a.Chan]Sem{},
	}
	tc(env, main, config)
}

func TestExecErr(t *testing.T) {
	f := func(env a.Environment, f a.Expname, expected error) {
		t.Helper()
		_, err := Exec(env, f)
		if err == nil {
			t.Fatalf("expecting error %q", expected)
		}
		if !errors.Is(err, expected) {
			t.Errorf("unexpected error; got %q; want %q", err, expected)
		}
	}

	f(make(a.Environment), "foo", ErrUndefinedProcess)
}
