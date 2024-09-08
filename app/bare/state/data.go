package state

import (
	"errors"

	"smecalculus/rolevod/lib/id"
)

type rootData struct {
	ID     string
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
	K    kind   `db:"kind"`
	ID   string `db:"id"`
	Name string `db:"name"`
}

type transition struct {
	FromID string `db:"from_id"`
	ToID   string `db:"to_id"`
	MsgID  string `db:"msg_id"`
	MsgKey string `db:"msg_key"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRef func(state) (Ref, error)
	// goverter:ignore K Name
	DataFromRef   func(Ref) state
	DataToRefs    func([]state) ([]Ref, error)
	DataFromRefs  func([]Ref) []state
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) []*rootData
)

func dataToRoot(dto rootData) Root {
	return dataToState(dto, dto.States[dto.ID])
}

func dataFromRoot(root Root) *rootData {
	if root == nil {
		return nil
	}
	dto := &rootData{
		ID:     root.get().String(),
		States: map[string]state{},
		Trs:    map[string][]transition{},
	}
	state := dataFromState(dto, root)
	dto.States[state.ID] = state
	return dto
}

func dataFromState(dto *rootData, root Root) state {
	switch st := root.(type) {
	case *TpRef:
		return state{K: refK, ID: st.ID.String()}
	case *One:
		return state{K: oneK, ID: st.ID.String()}
	case *Tensor:
		from := state{K: tensorK, ID: st.ID.String()}
		msg := dataFromState(dto, st.S)
		to := dataFromState(dto, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case *Lolli:
		from := state{K: lolliK, ID: st.ID.String()}
		msg := dataFromState(dto, st.S)
		to := dataFromState(dto, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case *With:
		from := state{K: withK, ID: st.ID.String()}
		for l, st := range st.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	case *Plus:
		from := state{K: plusK, ID: st.ID.String()}
		for l, st := range st.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	default:
		panic(ErrUnexpectedSt)
	}
}

func dataToState(root rootData, state state) Root {
	id, err := id.String[ID](state.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch state.K {
	case refK:
		return &TpRef{ID: id}
	case oneK:
		return &One{ID: id}
	case tensorK:
		tr := root.Trs[state.ID][0]
		return &Tensor{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case lolliK:
		tr := root.Trs[state.ID][0]
		return &Lolli{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case withK:
		st := &With{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	case plusK:
		st := &Plus{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	default:
		panic(ErrUnexpectedSt)
	}
}

var (
	errInvalidID = errors.New("invalid state id")
)
