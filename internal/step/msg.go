package step

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
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

type ProcRootMsg struct {
	ID   string   `json:"id"`
	PID  string   `json:"pid"`
	Term *TermMsg `json:"term"`
}

type TermKind string

const (
	Close = TermKind("close")
	Wait  = TermKind("wait")
	Send  = TermKind("send")
	Recv  = TermKind("recv")
	Lab   = TermKind("lab")
	Case  = TermKind("case")
	CTA   = TermKind("cta")
	Link  = TermKind("link")
	Spawn = TermKind("spawn")
	Fwd   = TermKind("fwd")
)

var termKindRequired = []validation.Rule{
	validation.Required,
	validation.In(Close, Wait, Send, Recv, Lab, Case, Spawn, Fwd, CTA),
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
	Fwd   *FwdMsg   `json:"fwd,omitempty"`
	CTA   *CTAMsg   `json:"cta,omitempty"`
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
		validation.Field(&mto.Fwd, validation.Required.When(mto.K == Fwd)),
		validation.Field(&mto.CTA, validation.Required.When(mto.K == CTA)),
	)
}

type CloseMsg struct {
	A ph.Msg `json:"a"`
}

func (mto CloseMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.A, validation.Required),
	)
}

type WaitMsg struct {
	X    ph.Msg  `json:"x"`
	Cont TermMsg `json:"cont"`
}

func (mto WaitMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.X, validation.Required),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type SendMsg struct {
	A ph.Msg `json:"a"`
	B ph.Msg `json:"b"`
}

func (mto SendMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.A, validation.Required),
		validation.Field(&mto.B, validation.Required),
	)
}

type RecvMsg struct {
	X    ph.Msg  `json:"x"`
	Y    ph.Msg  `json:"y"`
	Cont TermMsg `json:"cont"`
}

func (mto RecvMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.X, validation.Required),
		validation.Field(&mto.Y, validation.Required),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type LabMsg struct {
	A     ph.Msg `json:"a"`
	Label string `json:"label"`
}

func (mto LabMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.A, validation.Required),
		validation.Field(&mto.Label, core.NameRequired...),
	)
}

type CaseMsg struct {
	X   ph.Msg      `json:"x"`
	Brs []BranchMsg `json:"branches"`
}

func (mto CaseMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.X, validation.Required),
		validation.Field(&mto.Brs,
			validation.Required,
			validation.Length(1, 10),
			validation.Each(validation.Required),
		),
	)
}

type BranchMsg struct {
	Label string  `json:"label"`
	Cont  TermMsg `json:"cont"`
}

func (mto BranchMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Label, core.NameRequired...),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type SpawnMsg struct {
	PE   ph.Msg   `json:"pe"`
	CEs  []string `json:"ces"`
	Cont TermMsg  `json:"cont"`
	Sig  string   `json:"sig_id"`
}

func (mto SpawnMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.PE, validation.Required),
		validation.Field(&mto.CEs, core.CtxOptional...),
		validation.Field(&mto.Cont, validation.Required),
		validation.Field(&mto.Sig, id.Required...),
	)
}

type FwdMsg struct {
	C ph.Msg `json:"c"`
	D ph.Msg `json:"d"`
}

func (mto FwdMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.C, validation.Required),
		validation.Field(&mto.D, validation.Required),
	)
}

type CTAMsg struct {
	AK  string `json:"access_key"`
	Sig string `json:"sig_id"`
}

func (mto CTAMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.AK, id.Required...),
		validation.Field(&mto.Sig, id.Required...),
	)
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/chnl:Msg.*
// goverter:extend MsgFromTerm
// goverter:extend MsgToTerm
// goverter:extend MsgFromTermNilable
// goverter:extend MsgToTermNilable
var (
	MsgFromProcRoot func(ProcRoot) ProcRootMsg
	MsgToProcRoot   func(ProcRootMsg) (ProcRoot, error)
)

func MsgFromTermNilable(t Term) *TermMsg {
	if t == nil {
		return nil
	}
	term := MsgFromTerm(t)
	return &term
}

func MsgFromTerm(t Term) TermMsg {
	switch term := t.(type) {
	case CloseSpec:
		return TermMsg{
			K: Close,
			Close: &CloseMsg{
				A: ph.MsgFromPH(term.A),
			},
		}
	case WaitSpec:
		return TermMsg{
			K: Wait,
			Wait: &WaitMsg{
				X:    ph.MsgFromPH(term.X),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case SendSpec:
		return TermMsg{
			K: Send,
			Send: &SendMsg{
				A: ph.MsgFromPH(term.A),
				B: ph.MsgFromPH(term.B),
			},
		}
	case RecvSpec:
		return TermMsg{
			K: Recv,
			Recv: &RecvMsg{
				X:    ph.MsgFromPH(term.X),
				Y:    ph.MsgFromPH(term.Y),
				Cont: MsgFromTerm(term.Cont),
			},
		}
	case LabSpec:
		return TermMsg{
			K: Lab,
			Lab: &LabMsg{
				A:     ph.MsgFromPH(term.A),
				Label: string(term.L),
			},
		}
	case CaseSpec:
		brs := []BranchMsg{}
		for l, t := range term.Conts {
			brs = append(brs, BranchMsg{Label: string(l), Cont: MsgFromTerm(t)})
		}
		return TermMsg{
			K: Case,
			Case: &CaseMsg{
				X:   ph.MsgFromPH(term.X),
				Brs: brs,
			},
		}
	case SpawnSpec:
		return TermMsg{
			K: Spawn,
			Spawn: &SpawnMsg{
				PE:   ph.MsgFromPH(term.PE),
				CEs:  id.ConvertToStrings(term.CEs),
				Cont: MsgFromTerm(term.Cont),
				Sig:  term.Sig.String(),
			},
		}
	case FwdSpec:
		return TermMsg{
			K: Fwd,
			Fwd: &FwdMsg{
				C: ph.MsgFromPH(term.C),
				D: ph.MsgFromPH(term.D),
			},
		}
	case CTASpec:
		return TermMsg{
			K: CTA,
			CTA: &CTAMsg{
				Sig: term.Sig.String(),
				AK:  ak.ConvertToString(term.AK),
			},
		}
	default:
		panic(ErrTermTypeUnexpected(term))
	}
}

func MsgToTermNilable(mto *TermMsg) (Term, error) {
	if mto == nil {
		return nil, nil
	}
	return MsgToTerm(*mto)
}

func MsgToTerm(mto TermMsg) (Term, error) {
	switch mto.K {
	case Close:
		a, err := ph.MsgToPH(mto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case Wait:
		x, err := ph.MsgToPH(mto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case Send:
		a, err := ph.MsgToPH(mto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := ph.MsgToPH(mto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case Recv:
		x, err := ph.MsgToPH(mto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := ph.MsgToPH(mto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	case Lab:
		a, err := ph.MsgToPH(mto.Lab.A)
		if err != nil {
			return nil, err
		}
		return LabSpec{A: a, L: core.Label(mto.Lab.Label)}, nil
	case Case:
		x, err := ph.MsgToPH(mto.Case.X)
		if err != nil {
			return nil, err
		}
		conts := make(map[core.Label]Term, len(mto.Case.Brs))
		for _, b := range mto.Case.Brs {
			cont, err := MsgToTerm(b.Cont)
			if err != nil {
				return nil, err
			}
			conts[core.Label(b.Label)] = cont
		}
		return CaseSpec{X: x, Conts: conts}, nil
	case Spawn:
		pe, err := ph.MsgToPH(mto.Spawn.PE)
		if err != nil {
			return nil, err
		}
		ces, err := id.ConvertFromStrings(mto.Spawn.CEs)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToTerm(mto.Spawn.Cont)
		if err != nil {
			return nil, err
		}
		sigID, err := id.ConvertFromString(mto.Spawn.Sig)
		if err != nil {
			return nil, err
		}
		return SpawnSpec{PE: pe, CEs: ces, Cont: cont, Sig: sigID}, nil
	case Fwd:
		c, err := ph.MsgToPH(mto.Fwd.C)
		if err != nil {
			return nil, err
		}
		d, err := ph.MsgToPH(mto.Fwd.D)
		if err != nil {
			return nil, err
		}
		return FwdSpec{C: c, D: d}, nil
	case CTA:
		key, err := ak.ConvertFromString(mto.CTA.AK)
		if err != nil {
			return nil, err
		}
		sigID, err := id.ConvertFromString(mto.CTA.Sig)
		if err != nil {
			return nil, err
		}
		return CTASpec{AK: key, Sig: sigID}, nil
	default:
		panic(ErrUnexpectedTermKind(mto.K))
	}
}

func ErrUnexpectedTermKind(k TermKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}

func ErrUnexpectedStepKind(k StepKind) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
