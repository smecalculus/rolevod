package state

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type kind int

const (
	one = iota
	recur
	tensor
	lolli
	with
	plus
)

type refData struct {
	ID string `json:"id"`
	K  kind   `json:"kind"`
}

type rootData struct {
	InitialID string
	States    map[string]state
	Trs       map[string][]transition
}

type state struct {
	ID   string `db:"id"`
	K    kind   `db:"kind"`
	Name string `db:"name"`
}

type state2 struct {
	ID    string `db:"id"`
	K     kind   `db:"kind"`
	OnID  string `db:"msg_id"`
	OnKey string `db:"msg_key"`
}

type transition struct {
	FromID string `db:"from_id"`
	OnID   string `db:"msg_id"`
	OnKey  string `db:"msg_key"`
	ToID   string `db:"to_id"`
}

type transition2 struct {
	FromID string         `db:"from_id"`
	FromK  kind           `db:"from_kind"`
	OnID   sql.NullString `db:"msg_id"`
	OnKey  sql.NullString `db:"msg_key"`
	ToID   sql.NullString `db:"to_id"`
	ToK    sql.NullInt32  `db:"to_kind"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRefs    func([]*refData) ([]Ref, error)
	DataFromRefs  func([]Ref) []*refData
	DataToRoots   func([]*rootData) ([]Root, error)
	DataFromRoots func([]Root) []*rootData
)

func JsonFromRef(ref Ref) (sql.NullString, error) {
	dto := dataFromRef(ref)
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
	var dto refData
	err := json.Unmarshal([]byte(jsn.String), &dto)
	if err != nil {
		return nil, err
	}
	return dataToRef(&dto)
}

func dataFromRef(ref Ref) *refData {
	if ref == nil {
		return nil
	}
	rid := ref.RID().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return &refData{K: one, ID: rid}
	case RecurRef, RecurRoot:
		return &refData{K: recur, ID: rid}
	case TensorRef, TensorRoot:
		return &refData{K: tensor, ID: rid}
	case LolliRef, LolliRoot:
		return &refData{K: lolli, ID: rid}
	case WithRef, WithRoot:
		return &refData{K: with, ID: rid}
	case PlusRef, PlusRoot:
		return &refData{K: plus, ID: rid}
	default:
		panic(ErrUnexpectedRef(ref))
	}
}

func dataToRef(dto *refData) (Ref, error) {
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

func dataToRoot(dto *rootData) Root {
	return dataToState(dto, dto.States[dto.InitialID])
}

func dataToRoot2(dtos []transition2, initialID string) Root {
	states := map[string]state{}
	trs := map[string][]transition{}
	for _, tr := range dtos {
		states[tr.FromID] = state{ID: tr.FromID, K: kind(tr.FromK)}
		if !tr.ToID.Valid {
			continue
		}
		trs[tr.FromID] = append(trs[tr.FromID],
			transition{
				FromID: tr.FromID,
				OnID:   tr.OnID.String,
				OnKey:  tr.OnKey.String,
				ToID:   tr.ToID.String,
			})
	}
	dto := &rootData{States: states, Trs: trs}
	return dataToState(dto, states[initialID])
}

func dataFromRoot(root Root) *rootData {
	if root == nil {
		return nil
	}
	dto := &rootData{
		InitialID: root.RID().String(),
		States:    map[string]state{},
		Trs:       map[string][]transition{},
	}
	state := dataFromState(dto, root)
	dto.States[state.ID] = state
	return dto
}

func dataFromState(dto *rootData, r Root) state {
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
		tr := transition{FromID: from.ID, ToID: to.ID, OnID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case LolliRoot:
		from := state{K: lolli, ID: stateID}
		msg := dataFromState(dto, root.X)
		to := dataFromState(dto, root.Z)
		tr := transition{FromID: from.ID, ToID: to.ID, OnID: msg.ID}
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case WithRoot:
		from := state{K: with, ID: stateID}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, OnKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	case PlusRoot:
		from := state{K: plus, ID: stateID}
		for l, st := range root.Choices {
			to := dataFromState(dto, st)
			tr := transition{FromID: from.ID, ToID: to.ID, OnKey: string(l)}
			dto.States[to.ID] = to
			dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		}
		return from
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func dataToState(dto *rootData, state state) Root {
	stateID, err := id.String[ID](state.ID)
	if err != nil {
		panic(errInvalidID(state.ID))
	}
	switch state.K {
	case one:
		return OneRoot{ID: stateID}
	case recur:
		tr := dto.Trs[state.ID][0]
		toID, err := id.String[ID](tr.ToID)
		if err != nil {
			panic(errInvalidID(tr.ToID))
		}
		return RecurRoot{ID: stateID, Name: state.Name, ToID: toID}
	case tensor:
		tr := dto.Trs[state.ID][0]
		return TensorRoot{
			ID: stateID,
			A:  dataToState(dto, dto.States[tr.OnID]),
			C:  dataToState(dto, dto.States[tr.ToID]),
		}
	case lolli:
		tr := dto.Trs[state.ID][0]
		return LolliRoot{
			ID: stateID,
			X:  dataToState(dto, dto.States[tr.OnID]),
			Z:  dataToState(dto, dto.States[tr.ToID]),
		}
	case with:
		st := WithRoot{ID: stateID}
		for _, tr := range dto.Trs[state.ID] {
			st.Choices[Label(tr.OnKey)] = dataToState(dto, dto.States[tr.ToID])
		}
		return st
	case plus:
		st := PlusRoot{ID: stateID}
		for _, tr := range dto.Trs[state.ID] {
			st.Choices[Label(tr.OnKey)] = dataToState(dto, dto.States[tr.ToID])
		}
		return st
	default:
		panic(errUnexpectedKind(state.K))
	}
}

func errInvalidID(id string) error {
	return fmt.Errorf("invalid state id %q", id)
}

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
