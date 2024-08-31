package dcl

import (
	"errors"
	"smecalculus/rolevod/lib/core"
)

type tpRootData struct {
	ID     string
	Name   string
	States map[string]state
	Trs    map[string][]transition
}

type kind int

const (
	oneK = iota
	refK
	tensorK
	lolliK
	withK
	plusK
)

type state struct {
	K    kind
	ID   string
	Name string
}

type transition struct {
	FromID string
	ToID   string
	MsgID  string
	Label  string
}

func dataFromTpRoot(root TpRoot) tpRootData {
	data := &tpRootData{
		ID:     core.ToString[AR](root.ID),
		Name:   root.Name,
		States: make(map[string]state),
		Trs:    map[string][]transition{},
	}
	state := dataFromStype(data, root.St)
	data.States[state.ID] = state
	return *data
}

func dataFromStype(data *tpRootData, stype Stype) state {
	switch st := stype.(type) {
	case One:
		return state{K: oneK, ID: core.ToString[AR](st.ID)}
	case TpRef:
		return state{K: refK, ID: core.ToString[AR](st.ID)}
	case Tensor:
		from := state{K: tensorK, ID: core.ToString[AR](st.ID)}
		msg := dataFromStype(data, st.S)
		to := dataFromStype(data, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		data.States[to.ID] = to
		data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		return from
	case Lolli:
		from := state{K: lolliK, ID: core.ToString[AR](st.ID)}
		msg := dataFromStype(data, st.S)
		to := dataFromStype(data, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		data.States[to.ID] = to
		data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		return from
	case With:
		from := state{K: withK, ID: core.ToString[AR](st.ID)}
		for l, st := range st.Chs {
			to := dataFromStype(data, st)
			tr := transition{FromID: from.ID, ToID: to.ID, Label: string(l)}
			data.States[to.ID] = to
			data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		}
		return from
	case Plus:
		from := state{K: plusK, ID: core.ToString[AR](st.ID)}
		for l, st := range st.Chs {
			to := dataFromStype(data, st)
			tr := transition{FromID: from.ID, ToID: to.ID, Label: string(l)}
			data.States[to.ID] = to
			data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		}
		return from
	default:
		panic(ErrUnexpectedSt)
	}
}

func dataToTpRoot(data tpRootData) TpRoot {
	id, err := core.FromString[AR](data.ID)
	if err != nil {
		panic(errInvalidID)
	}
	return TpRoot{
		ID:   id,
		Name: data.Name,
		St:   dataToStype(data, data.States[data.ID]),
	}
}

func dataToStype(data tpRootData, from state) Stype {
	id, err := core.FromString[AR](from.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch from.K {
	case oneK:
		return One{ID: id}
	case refK:
		return TpRef{ID: id, Name: from.Name}
	case tensorK:
		tr := data.Trs[from.ID][0]
		return Tensor{
			ID: id,
			S:  dataToStype(data, data.States[tr.MsgID]),
			T:  dataToStype(data, data.States[tr.ToID]),
		}
	case lolliK:
		tr := data.Trs[from.ID][0]
		return Lolli{
			ID: id,
			S:  dataToStype(data, data.States[tr.MsgID]),
			T:  dataToStype(data, data.States[tr.ToID]),
		}
	case withK:
		st := With{ID: id}
		for _, tr := range data.Trs[from.ID] {
			st.Chs[Label(tr.Label)] = dataToStype(data, data.States[tr.ToID])
		}
		return st
	case plusK:
		st := Plus{ID: id}
		for _, tr := range data.Trs[from.ID] {
			st.Chs[Label(tr.Label)] = dataToStype(data, data.States[tr.ToID])
		}
		return st
	default:
		panic(ErrUnexpectedSt)
	}
}

var (
	errInvalidID         = errors.New("invalid state id")
	errUnexpectedTrsSize = errors.New("unexpected transitions size")
)
