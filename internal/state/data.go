package state

import (
	"database/sql"
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type kind int

const (
	unknown = iota
	one
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

type choiceData struct {
	Pat  string `json:"on"`
	ToID string `json:"to"`
}

type rootData struct {
	ID     string
	States map[string]state
}

type state struct {
	ID      string         `db:"id"`
	K       kind           `db:"kind"`
	FromID  sql.NullString `db:"from_id"`
	OnRef   *RefData       `db:"on_ref"`
	ToID    sql.NullString `db:"to_id"`
	Choices []choiceData   `db:"choices"`
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

func DataFromRef(ref Ref) *RefData {
	if ref == nil {
		return nil
	}
	rid := ref.RID().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return &RefData{K: one, ID: rid}
	case MenRef, MenRoot:
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
	rid, err := id.StringToID(dto.ID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case one:
		return OneRef{rid}, nil
	case recur:
		return MenRef{rid}, nil
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
	case MenRoot:
		dto.States[stID] = state{ID: stID, K: recur, FromID: fromID}
		return stID, nil
	case TensorRoot:
		onRef := DataFromRef(root.B)
		// onID, err := dataFromState(dto, root.B, "")
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
		// choices := make([]choiceData, len(root.Choices))
		var choices []choiceData
		for pat, toSt := range root.Choices {
			toID, err := dataFromState(dto, toSt, stID)
			if err != nil {
				return "", err
			}
			choices = append(choices, choiceData{string(pat), toID})
		}
		st := state{
			ID:      stID,
			K:       plus,
			FromID:  fromID,
			Choices: choices,
		}
		dto.States[stID] = st
		return stID, nil
	case WithRoot:
		// choices := make([]choiceData, len(root.Choices))
		var choices []choiceData
		for pat, toSt := range root.Choices {
			toID, err := dataFromState(dto, toSt, stID)
			if err != nil {
				return "", err
			}
			choices = append(choices, choiceData{string(pat), toID})
		}
		st := state{
			ID:      stID,
			K:       with,
			FromID:  fromID,
			Choices: choices,
		}
		dto.States[stID] = st
		return stID, nil
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func dataToState(states map[string]state, st state) (Root, error) {
	stID, err := id.StringToID(st.ID)
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
		chs := make(map[Label]Root, len(dto.Choices))
		for _, choice := range dto.Choices {
			ch, err := dataToState(states, states[choice.ToID])
			if err != nil {
				return nil, err
			}
			chs[Label(choice.Pat)] = ch
		}
		return PlusRoot{ID: stID, Choices: chs}, nil
	case with:
		dto := states[st.ID]
		chs := make(map[Label]Root, len(dto.Choices))
		for _, choice := range dto.Choices {
			ch, err := dataToState(states, states[choice.ToID])
			if err != nil {
				return nil, err
			}
			chs[Label(choice.Pat)] = ch
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
