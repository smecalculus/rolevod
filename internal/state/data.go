package state

import (
	"database/sql"
	"fmt"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type kind int

const (
	nonstate = iota
	one
	link
	tensor
	lolli
	plus
	with
)

type RefData struct {
	ID string `db:"id" json:"id"`
	K  kind   `db:"kind" json:"kind"`
}

type rootData struct {
	ID     string
	States map[string]state
}

type rootData2 struct {
	ID     string
	States []state
}

type state struct {
	ID      string         `db:"id"`
	K       kind           `db:"kind"`
	FromID  sql.NullString `db:"from_id"`
	FQN     sql.NullString `db:"fqn"`
	Pair    *pairData      `db:"pair"`
	Choices []choiceData   `db:"choices"`
}

type pairData struct {
	OnID string `json:"on"`
	ToID string `json:"to"`
}

type choiceData struct {
	OnPat string `json:"on"`
	ToID  string `json:"to"`
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
	case LinkRef, LinkRoot:
		return &RefData{K: link, ID: rid}
	case TensorRef, TensorRoot:
		return &RefData{K: tensor, ID: rid}
	case LolliRef, LolliRoot:
		return &RefData{K: lolli, ID: rid}
	case PlusRef, PlusRoot:
		return &RefData{K: plus, ID: rid}
	case WithRef, WithRoot:
		return &RefData{K: with, ID: rid}
	default:
		panic(ErrUnexpectedRefType(ref))
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
	case link:
		return LinkRef{rid}, nil
	case tensor:
		return TensorRef{rid}, nil
	case lolli:
		return LolliRef{rid}, nil
	case plus:
		return PlusRef{rid}, nil
	case with:
		return WithRef{rid}, nil
	default:
		panic(errUnexpectedKind(dto.K))
	}
}

func dataToRoot(dto *rootData) (Root, error) {
	return statesToRoot(dto.States, dto.States[dto.ID])
}

func dataToRoot2(dto *rootData2) (Root, error) {
	states := map[string]state{}
	for _, dto := range dto.States {
		states[dto.ID] = dto
	}
	return statesToRoot(states, states[dto.ID])
}

func dataFromRoot(root Root) *rootData {
	if root == nil {
		return nil
	}
	dto := &rootData{
		ID:     root.RID().String(),
		States: map[string]state{},
	}
	statesFromRoot(dto, root, "")
	return dto
}

func statesToRoot(states map[string]state, st state) (Root, error) {
	stID, err := id.StringToID(st.ID)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case one:
		return OneRoot{ID: stID}, nil
	case link:
		return LinkRoot{ID: stID, FQN: sym.StringToSym(st.FQN.String)}, nil
	case tensor:
		a, err := statesToRoot(states, states[st.Pair.OnID])
		if err != nil {
			return nil, err
		}
		c, err := statesToRoot(states, states[st.Pair.ToID])
		if err != nil {
			return nil, err
		}
		return TensorRoot{ID: stID, B: a, C: c}, nil
	case lolli:
		x, err := statesToRoot(states, states[st.Pair.OnID])
		if err != nil {
			return nil, err
		}
		z, err := statesToRoot(states, states[st.Pair.ToID])
		if err != nil {
			return nil, err
		}
		return LolliRoot{ID: stID, Y: x, Z: z}, nil
	case plus:
		dto := states[st.ID]
		choices := make(map[core.Label]Root, len(dto.Choices))
		for _, ch := range dto.Choices {
			choice, err := statesToRoot(states, states[ch.ToID])
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.OnPat)] = choice
		}
		return PlusRoot{ID: stID, Choices: choices}, nil
	case with:
		dto := states[st.ID]
		choices := make(map[core.Label]Root, len(dto.Choices))
		for _, ch := range dto.Choices {
			choice, err := statesToRoot(states, states[ch.ToID])
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.OnPat)] = choice
		}
		return WithRoot{ID: stID, Choices: choices}, nil
	default:
		panic(errUnexpectedKind(st.K))
	}
}

func statesFromRoot(dto *rootData, r Root, from string) (string, error) {
	var fromID sql.NullString
	if len(from) > 0 {
		fromID = sql.NullString{String: from, Valid: true}
	}
	stID := r.RID().String()
	switch root := r.(type) {
	case OneRoot:
		dto.States[stID] = state{ID: stID, K: one, FromID: fromID}
		return stID, nil
	case LinkRoot:
		fqn := sql.NullString{String: sym.StringFromSym(root.FQN), Valid: true}
		dto.States[stID] = state{ID: stID, K: link, FromID: fromID, FQN: fqn}
		return stID, nil
	case TensorRoot:
		onID, err := statesFromRoot(dto, root.B, stID)
		if err != nil {
			return "", err
		}
		toID, err := statesFromRoot(dto, root.C, stID)
		if err != nil {
			return "", err
		}
		toSt := state{
			ID:     stID,
			K:      tensor,
			FromID: fromID,
			Pair:   &pairData{OnID: onID, ToID: toID},
		}
		dto.States[stID] = toSt
		return stID, nil
	case LolliRoot:
		onID, err := statesFromRoot(dto, root.Y, stID)
		if err != nil {
			return "", err
		}
		toID, err := statesFromRoot(dto, root.Z, stID)
		if err != nil {
			return "", err
		}
		st := state{
			ID:     stID,
			K:      lolli,
			FromID: fromID,
			Pair:   &pairData{OnID: onID, ToID: toID},
		}
		dto.States[stID] = st
		return stID, nil
	case PlusRoot:
		// choices := make([]choiceData, len(root.Choices))
		var choices []choiceData
		for pat, toSt := range root.Choices {
			toID, err := statesFromRoot(dto, toSt, stID)
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
			toID, err := statesFromRoot(dto, toSt, stID)
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
		panic(ErrRootTypeUnexpected(r))
	}
}

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
