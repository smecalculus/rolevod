package step

import (
	"fmt"

	valid "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

type StepKind string

const (
	Proc = StepKind("proc")
	Msg  = StepKind("msg")
	Srv  = StepKind("srv")
)

type RootMsg struct {
	ID string   `json:"id"`
	K  StepKind `json:"kind"`
}

type TermKind string

const (
	Fwd   = TermKind("fwd")
	Spawn = TermKind("spawn")
	Close = TermKind("close")
	Wait  = TermKind("wait")
	Rec   = TermKind("ref")
	Send  = TermKind("send")
	Recv  = TermKind("recv")
	Lab   = TermKind("lab")
	Case  = TermKind("case")
)

type TermMsg struct {
	K     TermKind  `json:"kind"`
	Close *CloseMsg `json:"close,omitempty"`
	Wait  *WaitMsg  `json:"wait,omitempty"`
	Send  *SendMsg  `json:"send,omitempty"`
	Recv  *RecvMsg  `json:"recv,omitempty"`
	Lab   *LabMsg   `json:"lab,omitempty"`
	Case  *CaseMsg  `json:"case,omitempty"`
}

func (mto *TermMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.K,
			valid.Required,
			valid.In(
				Fwd, Spawn, Rec,
				Close, Wait,
				Send, Recv,
				Lab, Case,
			),
		),
		valid.Field(&mto.Close, valid.Required.When(mto.K == Close)),
		valid.Field(&mto.Wait, valid.Required.When(mto.K == Wait)),
		valid.Field(&mto.Send, valid.Required.When(mto.K == Send)),
		valid.Field(&mto.Recv, valid.Required.When(mto.K == Recv)),
		valid.Field(&mto.Lab, valid.Required.When(mto.K == Lab)),
		valid.Field(&mto.Case, valid.Required.When(mto.K == Case)),
	)
}

type CloseMsg struct {
	A string `json:"a"`
}

func (mto *CloseMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.A, valid.Required),
	)
}

type WaitMsg struct {
	X    string   `json:"x"`
	Cont *TermMsg `json:"cont"`
}

func (mto *WaitMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.X, valid.Required),
		valid.Field(&mto.Cont, valid.Required),
	)
}

type SendMsg struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (mto *SendMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.A, valid.Required),
		valid.Field(&mto.B, valid.Required),
	)
}

type RecvMsg struct {
	X    string   `json:"x"`
	Y    string   `json:"y"`
	Cont *TermMsg `json:"cont"`
}

func (mto *RecvMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.X, valid.Required),
		valid.Field(&mto.Y, valid.Required),
		valid.Field(&mto.Cont, valid.Required),
	)
}

type LabMsg struct {
	C     string `json:"c"`
	Label string `json:"label"`
}

func (mto *LabMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.C, valid.Required, valid.Length(20, 20)),
		valid.Field(&mto.Label, valid.Required, valid.Length(1, 64)),
	)
}

type CaseMsg struct {
	Z     string              `json:"z"`
	Conts map[string]*TermMsg `json:"conts"`
}

func (mto *CaseMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.Z, valid.Required, valid.Length(20, 20)),
		valid.Field(&mto.Conts, valid.Required, valid.Length(1, 10)),
	)
}

func MsgFromTerm(t Term) *TermMsg {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		return &TermMsg{
			K: Close,
			Close: &CloseMsg{
				A: term.A.String(),
			},
		}
	case WaitSpec:
		return &TermMsg{
			K: Wait,
			Wait: &WaitMsg{
				X:    term.X.String(),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case SendSpec:
		return &TermMsg{
			K: Send,
			Send: &SendMsg{
				A: term.A.String(),
				B: term.B.String(),
			},
		}
	case RecvSpec:
		return &TermMsg{
			K: Recv,
			Recv: &RecvMsg{
				X:    term.X.String(),
				Y:    term.Y.String(),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case LabSpec:
		return &TermMsg{
			K: Lab,
			Lab: &LabMsg{
				C:     term.C.String(),
				Label: string(term.L),
			},
		}
	case CaseSpec:
		conts := make(map[string]*TermMsg, len(term.Branches))
		for l, t := range term.Branches {
			conts[string(l)] = MsgFromTerm(t)
		}
		return &TermMsg{
			K: Case,
			Case: &CaseMsg{
				Z:     term.Z.String(),
				Conts: conts,
			},
		}
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func MsgToTerm(mto *TermMsg) (Term, error) {
	if mto == nil {
		return nil, nil
	}
	switch mto.K {
	case Close:
		a, err := id.StringToID(mto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case Wait:
		x, err := id.StringToID(mto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case Send:
		a, err := id.StringToID(mto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := id.StringToID(mto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case Recv:
		x, err := id.StringToID(mto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := id.StringToID(mto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	case Lab:
		c, err := id.StringToID(mto.Lab.C)
		if err != nil {
			return nil, err
		}
		return LabSpec{C: c, L: state.Label(mto.Lab.Label)}, nil
	case Case:
		z, err := id.StringToID(mto.Case.Z)
		if err != nil {
			return nil, err
		}
		branches := make(map[state.Label]Term, len(mto.Case.Conts))
		for l, t := range mto.Case.Conts {
			branch, err := MsgToTerm(t)
			if err != nil {
				return nil, err
			}
			branches[state.Label(l)] = branch
		}
		return CaseSpec{Z: z, Branches: branches}, nil
	default:
		panic(ErrUnexpectedTermKind(mto.K))
	}
}

func ErrUnexpectedTermKind(k TermKind) error {
	return fmt.Errorf("unexpected term kind %v", k)
}

func ErrUnexpectedStepKind(k StepKind) error {
	return fmt.Errorf("unexpected step kind %v", k)
}
