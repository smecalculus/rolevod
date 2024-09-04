package dcl

import (
	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/core"
)

type TpSpecMsg struct {
	Name string   `json:"name"`
	St   StypeMsg `json:"st"`
}

type ExpSpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type TpTeaserMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `json:"name"`
}

type TpRootMsg struct {
	ID   string   `param:"id" json:"id"`
	Name string   `json:"name"`
	St   StypeMsg `json:"st"`
}

type StypeMsg struct {
	K    Kind        `json:"kind"`
	ID   string      `json:"id"`
	Name string      `json:"name,omitempty"`
	M    *StypeMsg   `json:"message,omitempty"`
	S    *StypeMsg   `json:"session,omitempty"`
	Chs  []ChoiceMsg `json:"choices,omitempty"`
}

type ChoiceMsg struct {
	L string   `json:"label"`
	S StypeMsg `json:"session"`
}

type ExpRootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Kind string

const (
	OneK    = Kind("one")
	RefK    = Kind("ref")
	TensorK = Kind("tensor")
	LolliK  = Kind("lolli")
	WithK   = Kind("with")
	PlusK   = Kind("plus")
)

type OneMsg struct {
	K  Kind   `json:"kind"`
	ID string `json:"id"`
}

type TpRefMsg struct {
	K    Kind   `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend To.*
// goverter:extend msg.*
var (
	// tp
	MsgToTpSpec     func(TpSpecMsg) (TpSpec, error)
	MsgFromTpSpec   func(TpSpec) TpSpecMsg
	MsgFromTpRoot   func(TpRoot) TpRootMsg
	MsgToTpRoot     func(TpRootMsg) (TpRoot, error)
	MsgFromTpRoots  func([]TpRoot) []TpRootMsg
	MsgToTpRoots    func([]TpRootMsg) ([]TpRoot, error)
	MsgFromTpTeaser func(TpTeaser) TpTeaserMsg
	MsgToTpTeaser   func(TpTeaserMsg) (TpTeaser, error)
	// exp
	MsgToExpSpec   func(ExpSpecMsg) ExpSpec
	MsgFromExpSpec func(ExpSpec) ExpSpecMsg
	// goverter:ignore Ctx Zc
	MsgToExpRoot    func(ExpRootMsg) (ExpRoot, error)
	MsgFromExpRoot  func(ExpRoot) ExpRootMsg
	MsgFromExpRoots func([]ExpRoot) []ExpRootMsg
)

func msgFromStype(stype Stype) StypeMsg {
	switch st := stype.(type) {
	case One:
		return StypeMsg{K: OneK, ID: st.ID.String()}
	case TpRef:
		return StypeMsg{K: RefK, ID: st.ID.String(), Name: st.Name}
	case Tensor:
		m := msgFromStype(st.S)
		s := msgFromStype(st.T)
		return StypeMsg{
			K:  TensorK,
			ID: st.ID.String(),
			M:  &m,
			S:  &s,
		}
	case Lolli:
		m := msgFromStype(st.S)
		s := msgFromStype(st.T)
		return StypeMsg{
			K:  LolliK,
			ID: st.ID.String(),
			M:  &m,
			S:  &s,
		}
	case With:
		sts := make([]ChoiceMsg, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = ChoiceMsg{L: string(l), S: msgFromStype(st.Chs[l])}
		}
		return StypeMsg{K: WithK, ID: st.ID.String(), Chs: sts}
	case Plus:
		sts := make([]ChoiceMsg, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = ChoiceMsg{L: string(l), S: msgFromStype(st.Chs[l])}
		}
		return StypeMsg{K: PlusK, ID: st.ID.String(), Chs: sts}
	default:
		panic(ErrUnexpectedSt)
	}
}

func msgToStype(msg StypeMsg) (Stype, error) {
	id, err := core.FromString[AR](msg.ID)
	if err != nil {
		return nil, err
	}
	switch msg.K {
	case OneK:
		return One{ID: id}, nil
	case RefK:
		return TpRef{ID: id, Name: msg.Name}, nil
	case TensorK:
		m, err := msgToStype(*msg.M)
		if err != nil {
			return nil, err
		}
		s, err := msgToStype(*msg.S)
		if err != nil {
			return nil, err
		}
		return Tensor{ID: id, S: m, T: s}, nil
	case LolliK:
		m, err := msgToStype(*msg.M)
		if err != nil {
			return nil, err
		}
		s, err := msgToStype(*msg.S)
		if err != nil {
			return nil, err
		}
		return Lolli{ID: id, S: m, T: s}, nil
	case WithK:
		sts := make(Choices, len(msg.Chs))
		for _, ch := range msg.Chs {
			st, err := msgToStype(ch.S)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.L)] = st
		}
		return With{ID: id, Chs: sts}, nil
	case PlusK:
		sts := make(Choices, len(msg.Chs))
		for _, ch := range msg.Chs {
			st, err := msgToStype(ch.S)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.L)] = st
		}
		return Plus{ID: id, Chs: sts}, nil
	default:
		panic(ErrUnexpectedSt)
	}
}
