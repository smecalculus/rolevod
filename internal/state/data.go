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
	ID     string
	States map[string]state
}

type state struct {
	ID     string         `db:"id"`
	K      kind           `db:"kind"`
	FromID sql.NullString `db:"from_id"`
	OnRef  *RefData       `db:"on_ref"`
	ToID   sql.NullString `db:"to_id"`
	ToIDs  [][2]string    `db:"to_ids"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
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
	rid, err := id.StringTo(dto.ID)
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

func dataToRoot(dtos []state, rootID string) (Root, error) {
	states := map[string]state{}
	for _, dto := range dtos {
		states[dto.ID] = dto
	}
	return dataToState(states, states[rootID])
}

func dataToRoot2(dto *rootData) (Root, error) {
	states := map[string]state{}
	for _, st := range dto.States {
		states[st.ID] = st
	}
	return dataToState(states, states[dto.ID])
}

func dataFromRoot(root Root) *rootData {
	if root == nil {
		return nil
	}
	dto := &rootData{
		ID:     root.RID().String(),
		States: map[string]state{},
	}
	dataFromState(dto, root, "")
	return dto
}

func dataFromState(dto *rootData, r Root, from string) (string, error) {
	var fromID sql.NullString
	if len(from) > 0 {
		fromID = sql.NullString{String: from, Valid: true}
	}
	stID := r.RID().String()
	switch root := r.(type) {
	case OneRoot:
		dto.States[stID] = state{ID: stID, K: one, FromID: fromID}
		return stID, nil
	case RecurRoot:
		dto.States[stID] = state{ID: stID, K: recur, FromID: fromID}
		return stID, nil
	case TensorRoot:
		onRef := DataFromRef(root.B)
		toID, err := dataFromState(dto, root.C, stID)
		if err != nil {
			return "", err
		}
		st := state{
			ID:     stID,
			K:      tensor,
			FromID: fromID,
			OnRef:  onRef,
			ToID:   sql.NullString{String: toID, Valid: true},
		}
		dto.States[stID] = st
		return stID, nil
	case LolliRoot:
		onRef := DataFromRef(root.Y)
		toID, err := dataFromState(dto, root.Z, stID)
		if err != nil {
			return "", err
		}
		st := state{
			ID:     stID,
			K:      lolli,
			FromID: fromID,
			OnRef:  onRef,
			ToID:   sql.NullString{String: toID, Valid: true},
		}
		dto.States[stID] = st
		return stID, nil
	case PlusRoot:
		toIDs := make([][2]string, len(root.Choices))
		for l, s := range root.Choices {
			toID, err := dataFromState(dto, s, stID)
			if err != nil {
				return "", err
			}
			toIDs = append(toIDs, [2]string{string(l), toID})
		}
		st := state{
			ID:     stID,
			K:      plus,
			FromID: fromID,
			ToIDs:  toIDs,
		}
		dto.States[stID] = st
		return stID, nil
	case WithRoot:
		toIDs := make([][2]string, len(root.Choices))
		for l, s := range root.Choices {
			toID, err := dataFromState(dto, s, stID)
			if err != nil {
				return "", err
			}
			toIDs = append(toIDs, [2]string{string(l), toID})
		}
		st := state{
			ID:     stID,
			K:      with,
			FromID: fromID,
			ToIDs:  toIDs,
		}
		dto.States[stID] = st
		return stID, nil
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func dataToState(states map[string]state, st state) (Root, error) {
	stID, err := id.StringTo(st.ID)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case one:
		return OneRoot{ID: stID}, nil
	case tensor:
		a, err := DataToRef(st.OnRef)
		if err != nil {
			return nil, err
		}
		c, err := dataToState(states, states[st.ToID.String])
		if err != nil {
			return nil, err
		}
		return TensorRoot{ID: stID, B: a, C: c}, nil
	case lolli:
		x, err := DataToRef(st.OnRef)
		if err != nil {
			return nil, err
		}
		z, err := dataToState(states, states[st.ToID.String])
		if err != nil {
			return nil, err
		}
		return LolliRoot{ID: stID, Y: x, Z: z}, nil
	case plus:
		dto := states[st.ID]
		chs := make(map[Label]Root, len(dto.ToIDs))
		for _, pair := range dto.ToIDs {
			ch, err := dataToState(states, states[pair[1]])
			if err != nil {
				return nil, err
			}
			chs[Label(pair[0])] = ch
		}
		return PlusRoot{ID: stID, Choices: chs}, nil
	case with:
		dto := states[st.ID]
		chs := make(map[Label]Root, len(dto.ToIDs))
		for _, pair := range dto.ToIDs {
			ch, err := dataToState(states, states[pair[1]])
			if err != nil {
				return nil, err
			}
			chs[Label(pair[0])] = ch
		}
		return WithRoot{ID: stID, Choices: chs}, nil
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
