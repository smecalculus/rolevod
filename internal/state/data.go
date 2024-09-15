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
	one = iota
	tpRef
	tensor
	lolli
	with
	plus
)

type RefData struct {
	ID string `json:"id"`
	K  kind   `json:"kind"`
}

type state struct {
	ID   string `db:"id"`
	K    kind   `db:"kind"`
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
// goverter:extend DataFromRef
// goverter:extend DataToRef
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
		ID:     root.rootID().String(),
		States: map[string]state{},
		Trs:    map[string][]transition{},
	}
	state := dataFromState(dto, root)
	dto.States[state.ID] = state
	return dto
}

func dataFromState(dto *rootData, r Root) state {
	switch root := r.(type) {
	case TpRefRoot:
		return state{K: tpRef, ID: root.ID.String()}
	case OneRoot:
		return state{K: one, ID: root.ID.String()}
	case TensorRoot:
		from := state{K: tensor, ID: root.ID.String()}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case LolliRoot:
		from := state{K: lolli, ID: root.ID.String()}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case WithRoot:
		from := state{K: with, ID: root.ID.String()}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	case PlusRoot:
		from := state{K: plus, ID: root.ID.String()}
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
	case one:
		return OneRoot{ID: id}
	case tpRef:
		return TpRefRoot{ID: id}
	case tensor:
		tr := root.Trs[state.ID][0]
		return TensorRoot{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case lolli:
		tr := root.Trs[state.ID][0]
		return LolliRoot{
			ID: id,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case with:
		st := WithRoot{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	case plus:
		st := PlusRoot{ID: id}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	default:
		panic(errUnexpectedKind(state.K))
	}
}

func DataFromRef(r Ref) *RefData {
	if r == nil {
		return nil
	}
	switch ref := r.(type) {
	case OneRef:
		return &RefData{K: one, ID: ref.RootID().String()}
	case TpRefRef:
		return &RefData{K: tpRef, ID: ref.RootID().String()}
	case TensorRef:
		return &RefData{K: tensor, ID: ref.RootID().String()}
	case LolliRef:
		return &RefData{K: lolli, ID: ref.RootID().String()}
	case WithRef:
		return &RefData{K: with, ID: ref.RootID().String()}
	case PlusRef:
		return &RefData{K: plus, ID: ref.RootID().String()}
	default:
		panic(ErrUnexpectedRef(r))
	}
}

func DataToRef(dto *RefData) (Ref, error) {
	if dto == nil {
		return nil, nil
	}
	id, err := id.String[ID](dto.ID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case one:
		return OneRef(id), nil
	case tpRef:
		return TpRefRef(id), nil
	case tensor:
		return TensorRef(id), nil
	case lolli:
		return LolliRef(id), nil
	case with:
		return WithRef(id), nil
	case plus:
		return PlusRef(id), nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func JsonFromRef(ref Ref) (string, error) {
	dto := DataFromRef(ref)
	dtoJson, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}
	return string(dtoJson), nil
}

func JsonToRef(dtoJson string) (Ref, error) {
	var dto RefData
	err := json.Unmarshal([]byte(dtoJson), &dto)
	if err != nil {
		return nil, err
	}
	return DataToRef(&dto)
}

var (
	errInvalidID = errors.New("invalid state id")
)

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
