package step

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
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

var stepKindRequired = []validation.Rule{
	validation.Required,
	validation.In(Proc, Msg, Srv),
}

type RootMsg struct {
	ID string   `json:"id"`
	K  StepKind `json:"kind"`
}

type TermKind string

const (
	Close = TermKind("close")
	Wait  = TermKind("wait")
	Rec   = TermKind("ref")
	Send  = TermKind("send")
	Recv  = TermKind("recv")
	Lab   = TermKind("lab")
	Case  = TermKind("case")
	Spawn = TermKind("spawn")
	Fwd   = TermKind("fwd")
)

var termKindRequired = []validation.Rule{
	validation.Required,
	validation.In(Close, Wait, Send, Recv, Lab, Case, Spawn),
}

type TermMsg struct {
	K     TermKind  `json:"kind"`
	Close *CloseMsg `json:"close,omitempty"`
	Wait  *WaitMsg  `json:"wait,omitempty"`
	Send  *SendMsg  `json:"send,omitempty"`
	Recv  *RecvMsg  `json:"recv,omitempty"`
	Lab   *LabMsg   `json:"lab,omitempty"`
	Case  *CaseMsg  `json:"case,omitempty"`
	Spawn *SpawnMsg `json:"spawn,omitempty"`
}

func (mto TermMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.K, termKindRequired...),
		validation.Field(&mto.Close, validation.Required.When(mto.K == Close)),
		validation.Field(&mto.Wait, validation.Required.When(mto.K == Wait)),
		validation.Field(&mto.Send, validation.Required.When(mto.K == Send)),
		validation.Field(&mto.Recv, validation.Required.When(mto.K == Recv)),
		validation.Field(&mto.Lab, validation.Required.When(mto.K == Lab)),
		validation.Field(&mto.Case, validation.Required.When(mto.K == Case)),
		validation.Field(&mto.Spawn, validation.Required.When(mto.K == Spawn)),
	)
}

type CloseMsg struct {
	A string `json:"a"`
}

func (mto CloseMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.A, id.Required...),
	)
}

type WaitMsg struct {
	X    string  `json:"x"`
	Cont TermMsg `json:"cont"`
}

func (mto WaitMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.X, id.Required...),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type SendMsg struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (mto SendMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.A, id.Required...),
		validation.Field(&mto.B, id.Required...),
	)
}

type RecvMsg struct {
	X    string  `json:"x"`
	Y    string  `json:"y"`
	Cont TermMsg `json:"cont"`
}

func (mto RecvMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.X, id.Required...),
		validation.Field(&mto.Y, id.Required...),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type LabMsg struct {
	C     string `json:"c"`
	Label string `json:"label"`
}

func (mto LabMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.C, id.Required...),
		validation.Field(&mto.Label, core.NameRequired...),
	)
}

type CaseMsg struct {
	Z     string             `json:"z"`
	Conts map[string]TermMsg `json:"conts"`
}

func (mto CaseMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Z, id.Required...),
		validation.Field(&mto.Conts,
			validation.Required,
			validation.Length(1, 10),
			validation.Each(validation.Required),
		),
	)
}

type SpawnMsg struct {
	DecID string        `json:"dec_id"`
	C     string        `json:"c"`
	Ctx   []chnl.RefMsg `json:"ctx"`
	Cont  TermMsg       `json:"cont"`
}

func (mto SpawnMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.DecID, id.Required...),
		validation.Field(&mto.C, id.Required...),
		validation.Field(&mto.Ctx, core.CtxOptional...),
		validation.Field(&mto.Cont, validation.Required),
	)
}

func MsgFromTerm(t Term) TermMsg {
	switch term := t.(type) {
	case CloseSpec:
		return TermMsg{
			K: Close,
			Close: &CloseMsg{
				A: term.A.String(),
			},
		}
	case WaitSpec:
		return TermMsg{
			K: Wait,
			Wait: &WaitMsg{
				X:    term.X.String(),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case SendSpec:
		return TermMsg{
			K: Send,
			Send: &SendMsg{
				A: term.A.String(),
				B: term.B.String(),
			},
		}
	case RecvSpec:
		return TermMsg{
			K: Recv,
			Recv: &RecvMsg{
				X:    term.X.String(),
				Y:    term.Y.String(),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case LabSpec:
		return TermMsg{
			K: Lab,
			Lab: &LabMsg{
				C:     term.C.String(),
				Label: string(term.L),
			},
		}
	case CaseSpec:
		conts := make(map[string]TermMsg, len(term.Conts))
		for l, t := range term.Conts {
			conts[string(l)] = MsgFromTerm(t)
		}
		return TermMsg{
			K: Case,
			Case: &CaseMsg{
				Z:     term.Z.String(),
				Conts: conts,
			},
		}
	case SpawnSpec:
		// ctx := make([]chnl.RefMsg, len(term.Ctx))
		var ctx []chnl.RefMsg
		for name, chID := range term.Ctx {
			ctx = append(ctx, chnl.RefMsg{ID: chID.String(), Name: name})
		}
		return TermMsg{
			K: Spawn,
			Spawn: &SpawnMsg{
				DecID: term.DecID.String(),
				C:     term.C.String(),
				Ctx:   ctx,
				Cont:  MsgFromTerm(term.Cont),
			},
		}
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func MsgToTerm(mto TermMsg) (Term, error) {
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
		conts := make(map[state.Label]Term, len(mto.Case.Conts))
		for l, t := range mto.Case.Conts {
			cont, err := MsgToTerm(t)
			if err != nil {
				return nil, err
			}
			conts[state.Label(l)] = cont
		}
		return CaseSpec{Z: z, Conts: conts}, nil
	case Spawn:
		decID, err := id.StringToID(mto.Spawn.DecID)
		if err != nil {
			return nil, err
		}
		c, err := id.StringToID(mto.Spawn.C)
		if err != nil {
			return nil, err
		}
		ctx := make(map[chnl.Name]chnl.ID, len(mto.Spawn.Ctx))
		for _, ref := range mto.Spawn.Ctx {
			refID, err := id.StringToID(ref.ID)
			if err != nil {
				return nil, err
			}
			ctx[ref.Name] = refID
		}
		cont, err := MsgToTerm(mto.Spawn.Cont)
		if err != nil {
			return nil, err
		}
		return SpawnSpec{DecID: decID, C: c, Ctx: ctx, Cont: cont}, nil
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
