package state

import (
	"fmt"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/id"
)

type SpecMsg struct {
	K      Kind     `json:"kind"`
	TpRef  *RecMsg  `json:"tpref,omitempty"`
	Tensor *ProdMsg `json:"tensor,omitempty"`
	Lolli  *ProdMsg `json:"lolli,omitempty"`
	With   *SumMsg  `json:"with,omitempty"`
	Plus   *SumMsg  `json:"plus,omitempty"`
}

type RecMsg struct {
	Name string `json:"name"`
	ToID string `json:"to_id"`
}

type ProdMsg struct {
	Value *SpecMsg `json:"value"`
	State *SpecMsg `json:"state"`
}

type SumMsg struct {
	Choices []ChoiceMsg `json:"choices"`
}

type ChoiceMsg struct {
	Label string   `json:"label"`
	State *SpecMsg `json:"state"`
}

type RefMsg struct {
	ID string `param:"id" json:"id"`
	K  Kind   `json:"kind"`
}

type Kind string

const (
	One    = Kind("one")
	TpRef  = Kind("ref")
	Tensor = Kind("tensor")
	Lolli  = Kind("lolli")
	With   = Kind("with")
	Plus   = Kind("plus")
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend Msg.*
var (
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromSpec(s Spec) *SpecMsg {
	if s == nil {
		return nil
	}
	switch spec := s.(type) {
	case OneSpec:
		return &SpecMsg{K: One}
	case RecSpec:
		return &SpecMsg{
			K:     TpRef,
			TpRef: &RecMsg{ToID: spec.ToID.String(), Name: spec.Name}}
	case TensorSpec:
		return &SpecMsg{
			K: Tensor,
			Tensor: &ProdMsg{
				Value: MsgFromSpec(spec.S),
				State: MsgFromSpec(spec.T),
			},
		}
	case LolliSpec:
		return &SpecMsg{
			K: Lolli,
			Tensor: &ProdMsg{
				Value: MsgFromSpec(spec.S),
				State: MsgFromSpec(spec.T),
			},
		}
	case WithSpec:
		choices := make([]ChoiceMsg, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), State: MsgFromSpec(spec.Choices[l])}
		}
		return &SpecMsg{K: With, With: &SumMsg{Choices: choices}}
	case PlusSpec:
		choices := make([]ChoiceMsg, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), State: MsgFromSpec(spec.Choices[l])}
		}
		return &SpecMsg{K: Plus, Plus: &SumMsg{Choices: choices}}
	default:
		panic(ErrUnexpectedSpec(s))
	}
}

func MsgToSpec(mto *SpecMsg) (Spec, error) {
	if mto == nil {
		return nil, nil
	}
	switch mto.K {
	case One:
		return OneSpec{}, nil
	case TpRef:
		id, err := id.String[ID](mto.TpRef.ToID)
		if err != nil {
			return nil, err
		}
		return RecSpec{ToID: id, Name: mto.TpRef.Name}, nil
	case Tensor:
		v, err := MsgToSpec(mto.Tensor.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Tensor.State)
		if err != nil {
			return nil, err
		}
		return TensorSpec{S: v, T: s}, nil
	case Lolli:
		v, err := MsgToSpec(mto.Lolli.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Lolli.State)
		if err != nil {
			return nil, err
		}
		return LolliSpec{S: v, T: s}, nil
	case With:
		choices := make(map[Label]Spec, len(mto.With.Choices))
		for _, ch := range mto.With.Choices {
			choice, err := MsgToSpec(ch.State)
			if err != nil {
				return nil, err
			}
			choices[Label(ch.Label)] = choice
		}
		return WithSpec{Choices: choices}, nil
	case Plus:
		choices := make(map[Label]Spec, len(mto.Plus.Choices))
		for _, ch := range mto.Plus.Choices {
			choice, err := MsgToSpec(ch.State)
			if err != nil {
				return nil, err
			}
			choices[Label(ch.Label)] = choice
		}
		return PlusSpec{Choices: choices}, nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func MsgFromRef(ref Ref) *RefMsg {
	if ref == nil {
		return nil
	}
	id := ref.RootID().String()
	switch ref.(type) {
	case OneRef:
		return &RefMsg{K: One, ID: id}
	case RecRef:
		return &RefMsg{K: TpRef, ID: id}
	case TensorRef:
		return &RefMsg{K: Tensor, ID: id}
	case LolliRef:
		return &RefMsg{K: Lolli, ID: id}
	case WithRef:
		return &RefMsg{K: With, ID: id}
	case PlusRef:
		return &RefMsg{K: Plus, ID: id}
	default:
		panic(ErrUnexpectedRef(ref))
	}
}

func MsgToRef(mto *RefMsg) (Ref, error) {
	if mto == nil {
		return nil, nil
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return nil, err
	}
	switch mto.K {
	case One:
		return OneRef(id), nil
	case TpRef:
		return RecRef(id), nil
	case Tensor:
		return TensorRef(id), nil
	case Lolli:
		return LolliRef(id), nil
	case With:
		return WithRef(id), nil
	case Plus:
		return PlusRef(id), nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func ErrUnexpectedKind(k Kind) error {
	return fmt.Errorf("unextected kind %q", k)
}
