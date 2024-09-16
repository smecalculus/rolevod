package state

import (
	"database/sql"
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
	rec
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
	stateID := r.rootID().String()
	switch root := r.(type) {
	case OneRoot:
		return state{K: one, ID: stateID}
	case RecRoot:
		from := state{K: rec, ID: stateID}
		to := dto.States[root.ToID.String()]
		tr := transition{FromID: from.ID, ToID: to.ID}
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case TensorRoot:
		from := state{K: tensor, ID: stateID}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case LolliRoot:
		from := state{K: lolli, ID: stateID}
		msg := dataFromState(dto, root.S)
		to := dataFromState(dto, root.T)
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

func dataToState(root rootData, state state) Root {
	stateID, err := id.String[ID](state.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch state.K {
	case one:
		return OneRoot{ID: stateID}
	case rec:
		tr := root.Trs[state.ID][0]
		toID, err := id.String[ID](tr.ToID)
		if err != nil {
			panic(errInvalidID)
		}
		return RecRoot{ID: stateID, Name: state.Name, ToID: toID}
	case tensor:
		tr := root.Trs[state.ID][0]
		return TensorRoot{
			ID: stateID,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
		}
	case lolli:
		tr := root.Trs[state.ID][0]
		return LolliRoot{
			ID: stateID,
			S:  dataToState(root, root.States[tr.MsgID]),
			T:  dataToState(root, root.States[tr.ToID]),
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
	rid := ref.RootID().String()
	switch ref.(type) {
	case OneRef:
		return &RefData{K: one, ID: rid}
	case RecRef:
		return &RefData{K: rec, ID: rid}
	case TensorRef:
		return &RefData{K: tensor, ID: rid}
	case LolliRef:
		return &RefData{K: lolli, ID: rid}
	case WithRef:
		return &RefData{K: with, ID: rid}
	case PlusRef:
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
	case rec:
		return RecRef(rid), nil
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
	null := sql.NullString{}
	dto := DataFromRef(ref)
	if dto == nil {
		return null, nil
	}
	jsn, err := json.Marshal(dto)
	if err != nil {
		return null, err
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
