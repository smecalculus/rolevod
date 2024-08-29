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

type KindMsg string

const (
	OneK    = KindMsg("one")
	RefK    = KindMsg("ref")
	TensorK = KindMsg("tensor")
	LolliK  = KindMsg("lolli")
	WithK   = KindMsg("with")
	PlusK   = KindMsg("plus")
)

type OneMsg struct {
	K KindMsg `json:"tag"`
}

type TpRefMsg struct {
	K    KindMsg   `json:"tag"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ProductMsg struct {
	K KindMsg     `json:"tag"`
	M StypeMsg `json:"message"`
	S StypeMsg `json:"state"`
}

type SumMsg struct {
	K   KindMsg                `json:"tag"`
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
		return OneMsg{OneK}
	case TpRef:
		return TpRefMsg{RefK, st.Name, core.ToString[AR](st.ID)}
	case Tensor:
		return ProductMsg{TensorK, msgFromStype(st.S), msgFromStype(st.T)}
	case Lolli:
		return ProductMsg{LolliK, msgFromStype(st.S), msgFromStype(st.T)}
	case With:
		choices := make(map[string]StypeMsg, len(st.Chs))
		for l, st := range st.Chs {
			choices[string(l)] = msgFromStype(st)
		}
		return SumMsg{WithK, choices}
	case Plus:
		choices := make(map[string]StypeMsg, len(st.Chs))
		for l, st := range st.Chs {
			choices[string(l)] = msgFromStype(st)
		}
		return SumMsg{PlusK, choices}
	case nil:
		return nil
	default:
		panic(ErrUnexpectedSt)
	}
}
