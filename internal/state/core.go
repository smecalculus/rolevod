package state

import (
	"fmt"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/sym"
)

type ID = id.ADT

type Spec interface {
	spec()
}

type OneSpec struct{}

func (OneSpec) spec() {}

// aka TpName
type LinkSpec struct {
	Role sym.ADT
}

func (LinkSpec) spec() {}

type TensorSpec struct {
	B Spec
	C Spec
}

func (TensorSpec) spec() {}

type LolliSpec struct {
	Y Spec
	Z Spec
}

func (LolliSpec) spec() {}

// aka Internal Choice
type PlusSpec struct {
	Choices map[core.Label]Spec
}

func (PlusSpec) spec() {}

// aka External Choice
type WithSpec struct {
	Choices map[core.Label]Spec
}

func (WithSpec) spec() {}

type UpSpec struct {
	A Spec
}

func (UpSpec) spec() {}

type DownSpec struct {
	A Spec
}

func (DownSpec) spec() {}

type Ref interface {
	id.Identifiable
}

type OneRef struct {
	ID ID
}

func (r OneRef) Ident() ID { return r.ID }

type LinkRef struct {
	ID ID
}

func (r LinkRef) Ident() ID { return r.ID }

type PlusRef struct {
	ID ID
}

func (r PlusRef) Ident() ID { return r.ID }

type WithRef struct {
	ID ID
}

func (r WithRef) Ident() ID { return r.ID }

type TensorRef struct {
	ID ID
}

func (r TensorRef) Ident() ID { return r.ID }

type LolliRef struct {
	ID ID
}

func (r LolliRef) Ident() ID { return r.ID }

type UpRef struct {
	ID ID
}

func (r UpRef) Ident() ID { return r.ID }

type DownRef struct {
	ID ID
}

func (r DownRef) Ident() ID { return r.ID }

type Polarizable interface {
	Pol() Polarity
}

// aka Stype
type Root interface {
	id.Identifiable
	Polarizable
}

type Prod interface {
	Next() ID
}

type Sum interface {
	Next(core.Label) ID
}

type OneRoot struct {
	ID ID
}

func (OneRoot) spec() {}

func (r OneRoot) Ident() ID { return r.ID }

func (r OneRoot) Pol() Polarity { return Pos }

// aka TpName
type LinkRoot struct {
	ID   ID
	Role sym.ADT
}

func (LinkRoot) spec() {}

func (r LinkRoot) Ident() ID { return r.ID }

func (r LinkRoot) Pol() Polarity { return Zero }

// aka Internal Choice
type PlusRoot struct {
	ID      ID
	Choices map[core.Label]Root
}

func (r PlusRoot) Ident() ID { return r.ID }

func (r PlusRoot) Next(l core.Label) ID { return r.Choices[l].Ident() }

func (r PlusRoot) Pol() Polarity { return Pos }

// aka External Choice
type WithRoot struct {
	ID      ID
	Choices map[core.Label]Root
}

func (r WithRoot) Ident() ID { return r.ID }

func (r WithRoot) Next(l core.Label) ID { return r.Choices[l].Ident() }

func (r WithRoot) Pol() Polarity { return Neg }

type TensorRoot struct {
	ID ID
	B  Root // value
	C  Root // cont
}

func (r TensorRoot) Ident() ID { return r.ID }

func (r TensorRoot) Next() ID { return r.C.Ident() }

func (r TensorRoot) Pol() Polarity { return Pos }

type LolliRoot struct {
	ID ID
	Y  Root // value
	Z  Root // cont
}

func (r LolliRoot) Ident() ID { return r.ID }

func (r LolliRoot) Next() ID { return r.Z.Ident() }

func (r LolliRoot) Pol() Polarity { return Neg }

type UpRoot struct {
	ID ID
	A  Root
}

func (UpRoot) spec() {}

func (r UpRoot) Ident() ID { return r.ID }

func (r UpRoot) Pol() Polarity { return Zero }

type DownRoot struct {
	ID ID
	A  Root
}

func (DownRoot) spec() {}

func (r DownRoot) Ident() ID { return r.ID }

func (r DownRoot) Pol() Polarity { return Zero }

type TombstoneRoot struct {
	ID ID
}

func (TombstoneRoot) spec() {}

func (r TombstoneRoot) Ident() ID { return r.ID }

func (r TombstoneRoot) Pol() Polarity { return Zero }

type Polarity int

const (
	Pos  = Polarity(+1)
	Zero = Polarity(0)
	Neg  = Polarity(-1)
)

type Context struct {
	Linear map[ph.ADT]Root
}

// Endpoint aka ChanTp
type EP struct {
	Z  ph.ADT
	St Root
}

type Repo interface {
	Insert(Root) error
	SelectAll() ([]Ref, error)
	SelectByID(ID) (Root, error)
	SelectByIDs([]ID) ([]Root, error)
	SelectEnv([]ID) (map[ID]Root, error)
}

func ConvertSpecToRoot(s Spec) Root {
	if s == nil {
		return nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return OneRoot{ID: id.New()}
	case LinkSpec:
		return LinkRoot{ID: id.New(), Role: spec.Role}
	case TensorSpec:
		return TensorRoot{
			ID: id.New(),
			B:  ConvertSpecToRoot(spec.B),
			C:  ConvertSpecToRoot(spec.C),
		}
	case LolliSpec:
		return LolliRoot{
			ID: id.New(),
			Y:  ConvertSpecToRoot(spec.Y),
			Z:  ConvertSpecToRoot(spec.Z),
		}
	case WithSpec:
		choices := make(map[core.Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return WithRoot{ID: id.New(), Choices: choices}
	case PlusSpec:
		choices := make(map[core.Label]Root, len(spec.Choices))
		for lab, st := range spec.Choices {
			choices[lab] = ConvertSpecToRoot(st)
		}
		return PlusRoot{ID: id.New(), Choices: choices}
	default:
		panic(ErrSpecTypeUnexpected(spec))
	}
}

func CheckRef(got, want ID) error {
	if got != want {
		return fmt.Errorf("state mismatch: want %+v, got %+v", want, got)
	}
	return nil
}

// aka eqtp
func CheckRoot(got, want Root) error {
	switch wantSt := want.(type) {
	case OneRoot:
		_, ok := got.(OneRoot)
		if !ok {
			return ErrRootTypeMismatch(got, want)
		}
		return nil
	case TensorRoot:
		gotSt, ok := got.(TensorRoot)
		if !ok {
			return ErrRootTypeMismatch(got, want)
		}
		err := CheckRoot(gotSt.B, wantSt.B)
		if err != nil {
			return err
		}
		return CheckRoot(gotSt.C, wantSt.C)
	case LolliRoot:
		gotSt, ok := got.(LolliRoot)
		if !ok {
			return ErrRootTypeMismatch(got, want)
		}
		err := CheckRoot(gotSt.Y, wantSt.Y)
		if err != nil {
			return err
		}
		return CheckRoot(gotSt.Z, wantSt.Z)
	case PlusRoot:
		gotSt, ok := got.(PlusRoot)
		if !ok {
			return ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := CheckRoot(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	case WithRoot:
		gotSt, ok := got.(WithRoot)
		if !ok {
			return ErrRootTypeMismatch(got, want)
		}
		if len(gotSt.Choices) != len(wantSt.Choices) {
			return fmt.Errorf("choices mismatch: want %v items, got %v items", len(wantSt.Choices), len(gotSt.Choices))
		}
		for wantLab, wantChoice := range wantSt.Choices {
			gotChoice, ok := gotSt.Choices[wantLab]
			if !ok {
				return fmt.Errorf("label mismatch: want %q, got nothing", wantLab)
			}
			err := CheckRoot(gotChoice, wantChoice)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		panic(ErrRootTypeUnexpected(want))
	}
}

func ErrSpecTypeUnexpected(got Spec) error {
	return fmt.Errorf("spec type unexpected: %T", got)
}

func ErrRefTypeUnexpected(got Ref) error {
	return fmt.Errorf("ref type unexpected: %T", got)
}

func ErrDoesNotExist(want ID) error {
	return fmt.Errorf("root doesn't exist: %v", want)
}

func ErrMissingInEnv(want ID) error {
	return fmt.Errorf("root missing in env: %v", want)
}

func ErrMissingInCfg(want ID) error {
	return fmt.Errorf("root missing in cfg: %v", want)
}

func ErrRootTypeUnexpected(got Root) error {
	return fmt.Errorf("root type unexpected: %T", got)
}

func ErrRootTypeMismatch(got, want Root) error {
	return fmt.Errorf("root type mismatch: want %T, got %T", want, got)
}
