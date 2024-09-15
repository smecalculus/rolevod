package step

import (
	"fmt"

	valid "github.com/go-ozzo/ozzo-validation"

	"smecalculus/rolevod/internal/chnl"
)

type RefMsg struct {
	ID string `json:"id" param:"id"`
}

type StepKind string

const (
	ProcK = StepKind("proc")
	MsgK  = StepKind("msg")
	SrvK  = StepKind("srv")
)

type RootMsg struct {
	ID string   `json:"id"`
	K  StepKind `json:"kind"`
}

type TermKind string

const (
	FwdK   = TermKind("fwd")
	SpawnK = TermKind("spawn")
	CloseK = TermKind("close")
	WaitK  = TermKind("wait")
	RefK   = TermKind("ref")
	SendK  = TermKind("send")
	RecvK  = TermKind("recv")
	LabK   = TermKind("lab")
	CaseK  = TermKind("case")
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
			valid.In(FwdK, SpawnK, CloseK, WaitK, RefK, SendK, RecvK, LabK, CaseK),
		),
	)
}

type CloseMsg struct {
	A *chnl.RefMsg `json:"a"`
}

func (mto *CloseMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.A, valid.Required),
	)
}

type WaitMsg struct {
	X    *chnl.RefMsg `json:"x"`
	Cont *TermMsg     `json:"cont"`
}

func (mto *WaitMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.X, valid.Required),
		valid.Field(&mto.Cont, valid.Required),
	)
}

type SendMsg struct {
	A *chnl.RefMsg `json:"a"`
	B *chnl.RefMsg `json:"b"`
}

func (mto *SendMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.A, valid.Required),
		valid.Field(&mto.B, valid.Required),
	)
}

type RecvMsg struct {
	X    *chnl.RefMsg `json:"x"`
	Y    *chnl.RefMsg `json:"y"`
	Cont *TermMsg     `json:"cont"`
}

func (mto *RecvMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.X, valid.Required),
		valid.Field(&mto.Y, valid.Required),
		valid.Field(&mto.Cont, valid.Required),
	)
}

type LabMsg struct {
	C     *chnl.RefMsg `json:"c"`
	Label string       `json:"label"`
}

func (mto *LabMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.C, valid.Required),
		valid.Field(&mto.Label, valid.Required),
	)
}

type CaseMsg struct {
	Z     *chnl.RefMsg        `json:"z"`
	Conts map[string]*TermMsg `json:"conts"`
}

func (mto *CaseMsg) Validate() error {
	return valid.ValidateStruct(mto,
		valid.Field(&mto.Z, valid.Required),
		valid.Field(&mto.Conts, valid.Required, valid.Length(1, 10)),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRef  func(Ref) RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
	MsgFromRefs func([]Ref) []RefMsg
)

func MsgFromTerm(t Term) *TermMsg {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case Close:
		return &TermMsg{
			K: CloseK,
			Close: &CloseMsg{
				A: chnl.MsgFromRef(term.A),
			},
		}
	case Wait:
		x := chnl.MsgFromRef(term.X)
		return &TermMsg{
			K: WaitK,
			Wait: &WaitMsg{
				X:    x,
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case Send:
		return &TermMsg{
			K: SendK,
			Send: &SendMsg{
				A: chnl.MsgFromRef(term.A),
				B: chnl.MsgFromRef(term.B),
			},
		}
	case Recv:
		return &TermMsg{
			K: RecvK,
			Recv: &RecvMsg{
				X:    chnl.MsgFromRef(term.X),
				Y:    chnl.MsgFromRef(term.Y),
				Cont: MsgFromTerm(term.Cont),
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
	case CloseK:
		a, err := chnl.MsgToRef(*mto.Close.A)
		if err != nil {
			return nil, err
		}
		return Close{A: a}, nil
	case WaitK:
		x, err := chnl.MsgToRef(*mto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return Wait{X: x, Cont: cont}, nil
	case SendK:
		a, err := chnl.MsgToRef(*mto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := chnl.MsgToRef(*mto.Send.B)
		if err != nil {
			return nil, err
		}
		return Send{A: a, B: b}, nil
	case RecvK:
		x, err := chnl.MsgToRef(*mto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := chnl.MsgToRef(*mto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return Recv{X: x, Y: y, Cont: cont}, nil
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
