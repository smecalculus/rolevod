package dcl

import "smecalculus/rolevod/lib/core"

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" json:"id"`
}

type TpRootMsg struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	St   StypeMsg `json:"st"`
}

type ExpRootMsg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StypeMsg interface {
	stype()
}

func (OneMsg) stype()     {}
func (TpNameMsg) stype()  {}
func (ProductMsg) stype() {}
func (SumMsg) stype()     {}

type Tag string

const (
	one    = Tag("one")
	ref    = Tag("ref")
	tensor = Tag("tensor")
	lolli  = Tag("lolli")
	with   = Tag("with")
	plus   = Tag("plus")
)

type OneMsg struct {
	T Tag `json:"tag"`
}

type TpNameMsg struct {
	T    Tag    `json:"tag"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ProductMsg struct {
	T Tag      `json:"tag"`
	M StypeMsg `json:"message"`
	S StypeMsg `json:"state"`
}

type SumMsg struct {
	T   Tag                `json:"tag"`
	Chs map[Label]StypeMsg `json:"choices"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend msg.*
var (
	MsgToTpSpec     func(SpecMsg) TpSpec
	MsgFromTpSpec   func(TpSpec) SpecMsg
	MsgToExpSpec    func(SpecMsg) ExpSpec
	MsgFromExpSpec  func(ExpSpec) SpecMsg
	MsgFromTpRoot   func(TpRoot) TpRootMsg
	MsgFromTpRoots  func([]TpRoot) []TpRootMsg
	MsgFromExpRoot  func(ExpRoot) ExpRootMsg
	MsgFromExpRoots func([]ExpRoot) []ExpRootMsg
)

func msgFromStype(stype Stype) StypeMsg {
	switch st := stype.(type) {
	case One:
		return OneMsg{one}
	case TpName:
		return TpNameMsg{ref, st.Name, core.ToString[AR](st.ID)}
	case Tensor:
		return ProductMsg{tensor, msgFromStype(st.S), msgFromStype(st.T)}
	case Lolli:
		return ProductMsg{lolli, msgFromStype(st.S), msgFromStype(st.T)}
	case With:
		choices := make(map[Label]StypeMsg, len(st.Choices))
		for k, v := range st.Choices {
			choices[k] = msgFromStype(v)
		}
		return SumMsg{with, choices}
	case Plus:
		choices := make(map[Label]StypeMsg, len(st.Choices))
		for k, v := range st.Choices {
			choices[k] = msgFromStype(v)
		}
		return SumMsg{plus, choices}
	case nil:
		return nil
	default:
		panic(ErrUnexpectedSt)
	}
}
