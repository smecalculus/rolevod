package dcl

import "smecalculus/rolevod/lib/core"

type SpecMsg struct {
	Name string `json:"name"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type TpRootMsg struct {
	ID   string   `param:"id" json:"id"`
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
func (TpRefMsg) stype()   {}
func (ProductMsg) stype() {}
func (SumMsg) stype()     {}

type Tag string

const (
	OneT    = Tag("one")
	RefT    = Tag("ref")
	TensorT = Tag("tensor")
	LolliT  = Tag("lolli")
	WithT   = Tag("with")
	PlusT   = Tag("plus")
)

type OneMsg struct {
	T Tag `json:"tag"`
}

type TpRefMsg struct {
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
	T   Tag                 `json:"tag"`
	Chs map[string]StypeMsg `json:"choices"`
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
		return OneMsg{OneT}
	case TpRef:
		return TpRefMsg{RefT, st.Name, core.ToString[AR](st.ID)}
	case Tensor:
		return ProductMsg{TensorT, msgFromStype(st.S), msgFromStype(st.T)}
	case Lolli:
		return ProductMsg{LolliT, msgFromStype(st.S), msgFromStype(st.T)}
	case With:
		choices := make(map[string]StypeMsg, len(st.Choices))
		for k, v := range st.Choices {
			choices[string(k)] = msgFromStype(v)
		}
		return SumMsg{WithT, choices}
	case Plus:
		choices := make(map[string]StypeMsg, len(st.Choices))
		for k, v := range st.Choices {
			choices[string(k)] = msgFromStype(v)
		}
		return SumMsg{PlusT, choices}
	case nil:
		return nil
	default:
		panic(ErrUnexpectedSt)
	}
}
