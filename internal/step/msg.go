package step

import (
	"fmt"
	"smecalculus/rolevod/internal/chnl"
)

type RefMsg struct {
	ID string `param:"id" json:"id"`
}

type RootMsg struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

type RootKind string

const (
	ProcK = RootKind("proc")
	MsgK  = RootKind("msg")
	SrvK  = RootKind("srv")
)

type TermMsg struct {
	Kind  string    `json:"kind"`
	Close *CloseMsg `json:"close,omitempty"`
	Wait  *WaitMsg  `json:"wait,omitempty"`
	Send  *SendMsg  `json:"send,omitempty"`
	Recv  *RecvMsg  `json:"recv,omitempty"`
	Lab   *LabMsg   `json:"lab,omitempty"`
	Case  *CaseMsg  `json:"case,omitempty"`
}

type CloseMsg struct {
	A *chnl.RefMsg `json:"a"`
}

type WaitMsg struct {
	X    *chnl.RefMsg `json:"x"`
	Cont *TermMsg     `json:"cont"`
}

type SendMsg struct {
	A *chnl.RefMsg `json:"a"`
	B *chnl.RefMsg `json:"b"`
}

type RecvMsg struct {
	X    *chnl.RefMsg `json:"x"`
	Y    *chnl.RefMsg `json:"y"`
	Cont *TermMsg     `json:"cont"`
}

type LabMsg struct {
	C     *chnl.RefMsg `json:"c"`
	Label string       `json:"label"`
}

type CaseMsg struct {
	X     *chnl.RefMsg        `json:"x"`
	Conts map[string]*TermMsg `json:"conts"`
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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
var (
	MsgFromRef  func(Ref) RefMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromTerm(t Term) *TermMsg {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case Close:
		return &TermMsg{
			Kind: string(CloseK),
			Close: &CloseMsg{
				A: chnl.MsgFromRef(term.A),
			},
		}
	case Wait:
		fmt.Printf("%#v", term)
		fmt.Printf("%#v", term.X)
		x := chnl.MsgFromRef(term.X)
		return &TermMsg{
			Kind: string(WaitK),
			Wait: &WaitMsg{
				X:    x,
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case Send:
		return &TermMsg{
			Kind: string(SendK),
			Send: &SendMsg{
				A: chnl.MsgFromRef(term.A),
				B: chnl.MsgFromRef(term.B),
			},
		}
	case Recv:
		return &TermMsg{
			Kind: string(RecvK),
			Recv: &RecvMsg{
				X:    chnl.MsgFromRef(term.X),
				Y:    chnl.MsgFromRef(term.Y),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	default:
		panic(ErrUnexpectedTerm)
	}
}

func MsgToTerm(mto *TermMsg) (Term, error) {
	if mto == nil {
		return nil, nil
	}
	switch mto.Kind {
	case string(CloseK):
		a, err := chnl.MsgToRef(mto.Close.A)
		if err != nil {
			return nil, err
		}
		return Close{A: a}, nil
	case string(WaitK):
		x, err := chnl.MsgToRef(mto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return Wait{X: x, Cont: cont}, nil
	case string(SendK):
		a, err := chnl.MsgToRef(mto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := chnl.MsgToRef(mto.Send.B)
		if err != nil {
			return nil, err
		}
		return Send{A: a, B: b}, nil
	case string(RecvK):
		x, err := chnl.MsgToRef(mto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := chnl.MsgToRef(mto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return Recv{X: x, Y: y, Cont: cont}, nil
	default:
		panic(ErrUnexpectedKind)
	}
}
