package state

import (
	"encoding/json"
	"errors"
	"fmt"

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

type RefData struct {
	K  kind   `db:"kind" json:"kind"`
	ID string `db:"id" json:"id"`
}

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
	DataToRefs    func([]*RefData) ([]Ref, error)
	DataFromRefs  func([]Ref) []*RefData
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
		ID:     root.getID().String(),
		States: map[string]state{},
		Trs:    map[string][]transition{},
	}
	state := dataFromState(dto, root)
	dto.States[state.ID] = state
	return dto
}

func dataFromState(dto *rootData, r Root) state {
	switch root := r.(type) {
	case TpRef:
		return state{K: refK, ID: root.ID.String()}
	case One:
		return state{K: oneK, ID: root.ID.String()}
	case Tensor:
		from := state{K: tensorK, ID: root.ID.String()}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case Lolli:
		from := state{K: lolliK, ID: root.ID.String()}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case With:
		from := state{K: withK, ID: root.ID.String()}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	case Plus:
		from := state{K: plusK, ID: root.ID.String()}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func dataToState(root rootData, state state) Root {
	id, err := id.String[ID](state.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch state.K {
	case refK:
		return TpRef{ID: id}
	case oneK:
		return One{ID: id}
	case tensorK:
		tr := root.Trs[state.ID][0]
		return Tensor{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case lolliK:
		tr := root.Trs[state.ID][0]
		return Lolli{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case withK:
		st := With{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	case plusK:
		st := Plus{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	default:
		panic(errUnexpectedKind(state.K))
	}
}

func dataFromRef(r Ref) *RefData {
	if r == nil {
		return nil
	}
	switch ref := r.(type) {
	case OneRef:
		return &RefData{K: oneK, ID: ref.ID().String()}
	case TpRefRef:
		return &RefData{K: refK, ID: ref.ID().String()}
	case TensorRef:
		return &RefData{K: tensorK, ID: ref.ID().String()}
	case LolliRef:
		return &RefData{K: lolliK, ID: ref.ID().String()}
	case WithRef:
		return &RefData{K: withK, ID: ref.ID().String()}
	case PlusRef:
		return &RefData{K: plusK, ID: ref.ID().String()}
	default:
		panic(ErrUnexpectedRef(r))
	}
}

func dataToRef(dto *RefData) (Ref, error) {
	if dto == nil {
		return nil, nil
	}
	id, err := id.String[ID](dto.ID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case oneK:
		return OneRef{ref{id}}, nil
	case refK:
		return TpRefRef{ref{id}}, nil
	case tensorK:
		return TensorRef{ref{id}}, nil
	case lolliK:
		return LolliRef{ref{id}}, nil
	case withK:
		return WithRef{ref{id}}, nil
	case plusK:
		return PlusRef{ref{id}}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func JsonFromRef(ref Ref) (string, error) {
	dto := dataFromRef(ref)
	str, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func JsonToRef(data string) (Ref, error) {
	var dto RefData
	err := json.Unmarshal([]byte(data), &dto)
	if err != nil {
		return nil, err
	}
	return dataToRef(&dto)
}

var (
	errInvalidID = errors.New("invalid state id")
)

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
