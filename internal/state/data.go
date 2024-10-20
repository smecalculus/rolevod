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
	Type    typeData       `db:"type"`
}

type typeData struct {
	Tensor *prodData `json:"tensor,omitempty"`
	Lolli  *prodData `json:"lolli,omitempty"`
	Plus   []sumData `json:"plus,omitempty"`
	With   []sumData `json:"with,omitempty"`
}

type pairData struct {
	OnID string `json:"on"`
	ToID string `json:"to"`
}

type prodData struct {
	Val  string `json:"on"`
	Cont string `json:"to"`
}

type sumData struct {
	Lab  string `json:"on"`
	Cont string `json:"to"`
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
	DataToRoots   func([]*rootData2) ([]Root, error)
	DataFromRoots func([]Root) []*rootData2
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
		panic(ErrRefTypeUnexpected(ref))
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

func dataToRoot(dto *rootData2) (Root, error) {
	states := make(map[string]state, len(dto.States))
	for _, dto := range dto.States {
		states[dto.ID] = dto
	}
	return statesToRoot2(states, states[dto.ID])
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

func dataFromRoot2(root Root) *rootData2 {
	if root == nil {
		return nil
	}
	dto := &rootData2{
		ID:     root.RID().String(),
		States: nil,
	}
	statesFromRoot2("", root, dto)
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
		return LinkRoot{ID: stID, Role: sym.StringToSym(st.FQN.String)}, nil
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

func statesToRoot2(states map[string]state, st state) (Root, error) {
	stID, err := id.StringToID(st.ID)
	if err != nil {
		return nil, err
	}
	switch st.K {
	case one:
		return OneRoot{ID: stID}, nil
	case tensor:
		a, err := statesToRoot2(states, states[st.Type.Tensor.Val])
		if err != nil {
			return nil, err
		}
		c, err := statesToRoot2(states, states[st.Type.Tensor.Cont])
		if err != nil {
			return nil, err
		}
		return TensorRoot{ID: stID, B: a, C: c}, nil
	case lolli:
		x, err := statesToRoot2(states, states[st.Type.Lolli.Val])
		if err != nil {
			return nil, err
		}
		z, err := statesToRoot2(states, states[st.Type.Lolli.Cont])
		if err != nil {
			return nil, err
		}
		return LolliRoot{ID: stID, Y: x, Z: z}, nil
	case plus:
		dto := states[st.ID]
		choices := make(map[core.Label]Root, len(dto.Type.Plus))
		for _, ch := range dto.Type.Plus {
			choice, err := statesToRoot2(states, states[ch.Cont])
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Lab)] = choice
		}
		return PlusRoot{ID: stID, Choices: choices}, nil
	case with:
		dto := states[st.ID]
		choices := make(map[core.Label]Root, len(dto.Type.With))
		for _, ch := range dto.Type.With {
			choice, err := statesToRoot2(states, states[ch.Cont])
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Lab)] = choice
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
		fqn := sql.NullString{String: sym.StringFromSym(root.Role), Valid: true}
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

func statesFromRoot2(from string, r Root, dto *rootData2) (string, error) {
	var fromID sql.NullString
	if len(from) > 0 {
		fromID = sql.NullString{String: from, Valid: true}
	}
	stID := r.RID().String()
	switch root := r.(type) {
	case OneRoot:
		st := state{ID: stID, K: one, FromID: fromID}
		dto.States = append(dto.States, st)
		return stID, nil
	case TensorRoot:
		val, err := statesFromRoot2(stID, root.B, dto)
		if err != nil {
			return "", err
		}
		cont, err := statesFromRoot2(stID, root.C, dto)
		if err != nil {
			return "", err
		}
		st := state{
			ID:     stID,
			K:      tensor,
			FromID: fromID,
			Type: typeData{
				Tensor: &prodData{val, cont},
			},
		}
		dto.States = append(dto.States, st)
		return stID, nil
	case LolliRoot:
		val, err := statesFromRoot2(stID, root.Y, dto)
		if err != nil {
			return "", err
		}
		cont, err := statesFromRoot2(stID, root.Z, dto)
		if err != nil {
			return "", err
		}
		st := state{
			ID:     stID,
			K:      lolli,
			FromID: fromID,
			Type: typeData{
				Lolli: &prodData{val, cont},
			},
		}
		dto.States = append(dto.States, st)
		return stID, nil
	case PlusRoot:
		var choices []sumData
		for label, choice := range root.Choices {
			cont, err := statesFromRoot2(stID, choice, dto)
			if err != nil {
				return "", err
			}
			choices = append(choices, sumData{string(label), cont})
		}
		st := state{
			ID:     stID,
			K:      plus,
			FromID: fromID,
			Type:   typeData{Plus: choices},
		}
		dto.States = append(dto.States, st)
		return stID, nil
	case WithRoot:
		var choices []sumData
		for label, choice := range root.Choices {
			cont, err := statesFromRoot2(stID, choice, dto)
			if err != nil {
				return "", err
			}
			choices = append(choices, sumData{string(label), cont})
		}
		st := state{
			ID:     stID,
			K:      with,
			FromID: fromID,
			Type:   typeData{With: choices},
		}
		dto.States = append(dto.States, st)
		return stID, nil
	default:
		panic(ErrRootTypeUnexpected(r))
	}
}

func errUnexpectedKind(k kind) error {
	return fmt.Errorf("unexpected kind %q", k)
}
