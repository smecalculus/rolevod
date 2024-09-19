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

type RefData struct {
	ID string `db:"id" json:"id"`
	K  kind   `db:"kind" json:"kind"`
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
// goverter:extend DataToRef
// goverter:extend DataFromRef
var (
	DataToRefs    func([]*RefData) ([]Ref, error)
	DataFromRefs  func([]Ref) []*RefData
	DataToRoots   func([]*rootData) ([]Root, error)
	DataFromRoots func([]Root) []*rootData
)

func JsonFromRef(ref Ref) (sql.NullString, error) {
	if ref == nil {
		return sql.NullString{}, nil
	}
	jsn, err := json.Marshal(DataFromRef(ref))
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
		return OneRef{rid}, nil
	case recur:
		return RecurRef{rid}, nil
	case tensor:
		return TensorRef{rid}, nil
	case lolli:
		return LolliRef{rid}, nil
	case with:
		return WithRef{rid}, nil
	case plus:
		return PlusRef{rid}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToRoot(dto *rootData) (Root, error) {
	return dataToState(dto, dto.States[dto.InitialID])
}

func dataToRoot2(dtos []transition2, initialID string) (Root, error) {
	states := map[string]state{}
	trs := map[string][]transition{}
	for _, tr := range dtos {
		states[tr.FromID] = state{ID: tr.FromID, K: kind(tr.FromK)}
		if !tr.ToID.Valid {
			continue
		}
		_, ok := trs[tr.FromID]
		if !ok {
			trs[tr.FromID] = []transition{}
		}
		trs[tr.FromID] = append(trs[tr.FromID],
			transition{
				FromID: tr.FromID,
				OnID:   tr.OnID.String,
				OnKey:  tr.OnKey.String,
				ToID:   tr.ToID.String,
			})
		states[tr.ToID.String] = state{ID: tr.ToID.String, K: kind(tr.ToK.Int32)}
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
		on := dataFromState(dto, root.A)
		to := dataFromState(dto, root.C)
		tr := transition{FromID: from.ID, ToID: to.ID, OnID: on.ID}
		dto.States[on.ID] = on
		dto.States[to.ID] = to
		dto.Trs[from.ID] = append(dto.Trs[from.ID], tr)
		return from
	case LolliRoot:
		from := state{K: lolli, ID: stateID}
		on := dataFromState(dto, root.X)
		to := dataFromState(dto, root.Z)
		tr := transition{FromID: from.ID, ToID: to.ID, OnID: on.ID}
		dto.States[on.ID] = on
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

func dataToState(dto *rootData, st state) (Root, error) {
	stID, err := id.String[ID](st.ID)
	if err != nil {
		panic(err)
	}
	switch st.K {
	case one:
		return OneRoot{ID: stID}, nil
	case recur:
		tr := dto.Trs[st.ID][0]
		toID, err := id.String[ID](tr.ToID)
		if err != nil {
			panic(errInvalidID(tr.ToID))
		}
		return RecurRoot{ID: stID, Name: st.Name, ToID: toID}, nil
	case tensor:
		tr := dto.Trs[st.ID][0]
		fmt.Printf("transition: %#v\n", tr)
		fmt.Printf("states: %#v\n", dto.States)
		a, err := dataToState(dto, dto.States[tr.OnID])
		if err != nil {
			return nil, err
		}
		c, err := dataToState(dto, dto.States[tr.ToID])
		if err != nil {
			return nil, err
		}
		return TensorRoot{ID: stID, A: a, C: c}, nil
	case lolli:
		tr := dto.Trs[st.ID][0]
		x, err := dataToState(dto, dto.States[tr.OnID])
		if err != nil {
			return nil, err
		}
		z, err := dataToState(dto, dto.States[tr.ToID])
		if err != nil {
			return nil, err
		}
		return LolliRoot{ID: stID, X: x, Z: z}, nil
	case with:
		state := WithRoot{ID: stID}
		for _, tr := range dto.Trs[st.ID] {
			ch, err := dataToState(dto, dto.States[tr.ToID])
			if err != nil {
				return nil, err
			}
			state.Choices[Label(tr.OnKey)] = ch
		}
		return state, nil
	case plus:
		state := PlusRoot{ID: stID}
		for _, tr := range dto.Trs[st.ID] {
			ch, err := dataToState(dto, dto.States[tr.ToID])
			if err != nil {
				return nil, err
			}
			state.Choices[Label(tr.OnKey)] = ch
		}
		return state, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func errInvalidID(id string) error {
	return fmt.Errorf("invalid state id %q", id)
}

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
