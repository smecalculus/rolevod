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
	K  KindMsg `json:"kind"`
	ID string  `json:"id"`
}

type TpRefMsg struct {
	K    KindMsg `json:"kind"`
	ID   string  `json:"id"`
	Name string  `json:"name"`
}

type ProductMsg struct {
	K  KindMsg  `json:"kind"`
	ID string   `json:"id"`
	M  StypeMsg `json:"message"`
	S  StypeMsg `json:"state"`
}

type SumMsg struct {
	K   KindMsg             `json:"kind"`
	ID  string              `json:"id"`
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
		return OneMsg{K: OneK, ID: core.ToString[AR](st.ID)}
	case TpRef:
		return TpRefMsg{K: RefK, ID: core.ToString[AR](st.ID), Name: st.Name}
	case Tensor:
		return ProductMsg{
			K:  TensorK,
			ID: core.ToString[AR](st.ID),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case Lolli:
		return ProductMsg{
			K:  LolliK,
			ID: core.ToString[AR](st.ID),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case With:
		sts := make(map[string]StypeMsg, len(st.Chs))
		for l, st := range st.Chs {
			sts[string(l)] = msgFromStype(st)
		}
		return SumMsg{K: WithK, ID: core.ToString[AR](st.ID), Chs: sts}
	case Plus:
		sts := make(map[string]StypeMsg, len(st.Chs))
		for l, st := range st.Chs {
			sts[string(l)] = msgFromStype(st)
		}
		return SumMsg{K: PlusK, ID: core.ToString[AR](st.ID), Chs: sts}
	case nil:
		return nil
	default:
		panic(ErrUnexpectedSt)
	}
}
