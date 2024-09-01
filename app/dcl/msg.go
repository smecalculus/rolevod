package dcl

import (
	"smecalculus/rolevod/lib/core"

	"golang.org/x/exp/maps"
)

type TpSpecMsg struct {
	Name string   `json:"name"`
	St   StypeRaw `json:"st"`
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

type TpRootRaw struct {
	ID   string   `param:"id" json:"id"`
	Name string   `json:"name"`
	St   StypeRaw `json:"st"`
}

type StypeRaw struct {
	K    Kind        `json:"kind"`
	ID   string      `json:"id"`
	Name string      `json:"name"`
	M    *StypeRaw   `json:"message"`
	S    *StypeRaw   `json:"session"`
	Chs  []ChoiceRaw `json:"choices"`
}

type ChoiceRaw struct {
	L string   `json:"label"`
	S StypeRaw `json:"session"`
}

type ExpRootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StypeMsg interface {
	stype()
}

func (OneMsg) stype()   {}
func (TpRefMsg) stype() {}
func (Product) stype()  {}
func (Sum) stype()      {}

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

type Product struct {
	K  Kind     `json:"kind"`
	ID string   `json:"id"`
	M  StypeMsg `json:"message"`
	S  StypeMsg `json:"session"`
}

type Sum struct {
	K   Kind     `json:"kind"`
	ID  string   `json:"id"`
	Chs []Choice `json:"choices"`
}

type Choice struct {
	L string   `json:"label"`
	S StypeMsg `json:"session"`
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
	MsgToTpRoot     func(TpRootRaw) (TpRoot, error)
	MsgFromTpRoots  func([]TpRoot) []TpRootMsg
	MsgToTpRoots    func([]TpRootRaw) ([]TpRoot, error)
	MsgFromTpTeaser func(TpTeaser) TpTeaserMsg
	MsgToTpTeaser   func(TpTeaserMsg) (TpTeaser, error)
	// exp
	MsgToExpSpec    func(ExpSpecMsg) ExpSpec
	MsgFromExpSpec  func(ExpSpec) ExpSpecMsg
	MsgToExpRoot    func(ExpRootMsg) (ExpRoot, error)
	MsgFromExpRoot  func(ExpRoot) ExpRootMsg
	MsgFromExpRoots func([]ExpRoot) []ExpRootMsg
)

func msgFromStype(stype Stype) StypeMsg {
	switch st := stype.(type) {
	case One:
		return OneMsg{K: OneK, ID: core.ToString[AR](st.ID)}
	case TpRef:
		return TpRefMsg{K: RefK, ID: core.ToString[AR](st.ID), Name: st.Name}
	case Tensor:
		return Product{
			K:  TensorK,
			ID: core.ToString[AR](st.ID),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case Lolli:
		return Product{
			K:  LolliK,
			ID: core.ToString[AR](st.ID),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case With:
		sts := make([]Choice, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = Choice{L: string(l), S: msgFromStype(st.Chs[l])}
		}
		return Sum{K: WithK, ID: core.ToString[AR](st.ID), Chs: sts}
	case Plus:
		sts := make([]Choice, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = Choice{L: string(l), S: msgFromStype(st.Chs[l])}
		}
		return Sum{K: PlusK, ID: core.ToString[AR](st.ID), Chs: sts}
	case nil:
		return nil
	default:
		panic(ErrUnexpectedSt)
	}
}

func msgFromStype1(stype Stype) StypeRaw {
	switch st := stype.(type) {
	case One:
		return StypeRaw{K: OneK, ID: core.ToString(st.ID)}
	case TpRef:
		return StypeRaw{K: RefK, ID: core.ToString(st.ID), Name: st.Name}
	case Tensor:
		m := msgFromStype1(st.S)
		s := msgFromStype1(st.T)
		return StypeRaw{
			K:  TensorK,
			ID: core.ToString(st.ID),
			M:  &m,
			S:  &s,
		}
	case Lolli:
		m := msgFromStype1(st.S)
		s := msgFromStype1(st.T)
		return StypeRaw{
			K:  LolliK,
			ID: core.ToString[AR](st.ID),
			M:  &m,
			S:  &s,
		}
	case With:
		sts := make([]ChoiceRaw, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = ChoiceRaw{L: string(l), S: msgFromStype1(st.Chs[l])}
		}
		return StypeRaw{K: WithK, ID: core.ToString(st.ID), Chs: sts}
	case Plus:
		sts := make([]ChoiceRaw, len(st.Chs))
		for i, l := range maps.Keys(st.Chs) {
			sts[i] = ChoiceRaw{L: string(l), S: msgFromStype1(st.Chs[l])}
		}
		return StypeRaw{K: PlusK, ID: core.ToString(st.ID), Chs: sts}
	default:
		panic(ErrUnexpectedSt)
	}
}

func msgToStype(msg StypeRaw) (Stype, error) {
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
