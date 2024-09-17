package state

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type RootData struct {
	ID     string
	States map[string]state
	Trs    map[string][]transition
}

type kind int

const (
	one = iota
	recur
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
	DataToRoots   func([]*RootData) ([]Root, error)
	DataFromRoots func([]Root) []*RootData
)

func dataToRoot(dto *RootData) Root {
	return dataToState(dto, dto.States[dto.ID])
}

func dataFromRoot(root Root) *RootData {
	if root == nil {
		return nil
	}
	dto := &RootData{
		ID:     root.RID().String(),
		States: map[string]state{},
		Trs:    map[string][]transition{},
	}
	state := dataFromState(dto, root)
	dto.States[state.ID] = state
	return dto
}

func dataFromState(dto *RootData, r Root) state {
	stateID := r.RID().String()
	switch root := r.(type) {
	case OneRoot:
		return state{K: one, ID: stateID}
	case RecurRoot:
		from := state{K: recur, ID: stateID}
		to := dto.States[root.ToID.String()]
		tr := transition{FromID: from.ID, ToID: to.ID}
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case TensorRoot:
		from := state{K: tensor, ID: stateID}
		msg := dataFromState(dto, root.A)
		to := dataFromState(dto, root.C)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case LolliRoot:
		from := state{K: lolli, ID: stateID}
		msg := dataFromState(dto, root.X)
		to := dataFromState(dto, root.Z)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case WithRoot:
		from := state{K: with, ID: stateID}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	case PlusRoot:
		from := state{K: plus, ID: stateID}
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

func dataToState(root *RootData, state state) Root {
	stateID, err := id.String[ID](state.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch state.K {
	case one:
		return OneRoot{ID: stateID}
	case recur:
		tr := root.Trs[state.ID][0]
		toID, err := id.String[ID](tr.ToID)
		if err != nil {
			panic(errInvalidID)
		}
		return RecurRoot{ID: stateID, Name: state.Name, ToID: toID}
	case tensor:
		tr := root.Trs[state.ID][0]
		return TensorRoot{
			ID: stateID,
			A:  dataToState(root, root.States[tr.MsgID]),
			C:  dataToState(root, root.States[tr.ToID]),
		}
	case lolli:
		tr := root.Trs[state.ID][0]
		return LolliRoot{
			ID: stateID,
			X:  dataToState(root, root.States[tr.MsgID]),
			Z:  dataToState(root, root.States[tr.ToID]),
		}
	case with:
		st := WithRoot{ID: stateID}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	case plus:
		st := PlusRoot{ID: stateID}
		for _, tr := range root.Trs[state.ID] {
			st.Choices[Label(tr.MsgKey)] = dataToState(root, root.States[tr.ToID])
		}
		return st
	default:
		panic(errUnexpectedKind(state.K))
	}
}

func DataFromRef(ref Ref) *RefData {
	if ref == nil {
		return nil
	}
	rid := ref.RID().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return &RefData{K: one, ID: rid}
	case RecurRef, RecurRoot:
		return &RefData{K: recur, ID: rid}
	case TensorRef, TensorRoot:
		return &RefData{K: tensor, ID: rid}
	case LolliRef, LolliRoot:
		return &RefData{K: lolli, ID: rid}
	case WithRef, WithRoot:
		return &RefData{K: with, ID: rid}
	case PlusRef, PlusRoot:
		return &RefData{K: plus, ID: rid}
	default:
		panic(ErrUnexpectedRef(ref))
	}
}

func DataToRef(dto *RefData) (Ref, error) {
	if dto == nil {
		return nil, nil
	}
	rid, err := id.String[ID](dto.ID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case one:
		return OneRef(rid), nil
	case recur:
		return RecurRef(rid), nil
	case tensor:
		return TensorRef(rid), nil
	case lolli:
		return LolliRef(rid), nil
	case with:
		return WithRef(rid), nil
	case plus:
		return PlusRef(rid), nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func JsonFromRef(ref Ref) (sql.NullString, error) {
	dto := DataFromRef(ref)
	if dto == nil {
		return sql.NullString{}, nil
	}
	jsn, err := json.Marshal(dto)
	if err != nil {
		return sql.NullString{}, err
	}
	return sql.NullString{String: string(jsn), Valid: true}, nil
}

func JsonToRef(jsn sql.NullString) (Ref, error) {
	if !jsn.Valid {
		return nil, nil
	}
	var dto RefData
	err := json.Unmarshal([]byte(jsn.String), &dto)
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
