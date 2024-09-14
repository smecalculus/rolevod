package state

import (
	"fmt"

	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/id"
)

type RefMsg struct {
	ID string `param:"id" json:"id"`
	K  Kind   `json:"kind"`
}

type RootMsg struct {
	ID      string      `json:"id"`
	K       Kind        `json:"kind"`
	Name    string      `json:"name,omitempty"`
	Value   *RootMsg    `json:"value,omitempty"`
	State   *RootMsg    `json:"state,omitempty"`
	Choices []ChoiceMsg `json:"choices,omitempty"`
}

type ChoiceMsg struct {
	Label string   `json:"label"`
	State *RootMsg `json:"state"`
}

type Kind string

const (
	OneK    = Kind("one")
	RefK    = Kind("ref")
	TensorK = Kind("tensor")
	LolliK  = Kind("lolli")
	WithK   = Kind("with")
	PlusK   = Kind("plus")
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend Msg.*
var (
	ToRefMsg    func(*RootMsg) *RefMsg
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromRef(ref Ref) *RefMsg {
	if ref == nil {
		return nil
	}
	id := ref.ID().String()
	switch ref.(type) {
	case OneRef:
		return &RefMsg{K: OneK, ID: id}
	case TpRefRef:
		return &RefMsg{K: RefK, ID: id}
	case TensorRef:
		return &RefMsg{K: TensorK, ID: id}
	case LolliRef:
		return &RefMsg{K: LolliK, ID: id}
	case WithRef:
		return &RefMsg{K: WithK, ID: id}
	case PlusRef:
		return &RefMsg{K: PlusK, ID: id}
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
	case RefK:
		return TpRefRef{ref{id}}, nil
	case OneK:
		return OneRef{ref{id}}, nil
	case TensorK:
		return TensorRef{ref{id}}, nil
	case LolliK:
		return LolliRef{ref{id}}, nil
	case WithK:
		return WithRef{ref{id}}, nil
	case PlusK:
		return PlusRef{ref{id}}, nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func MsgFromRoot(r Root) *RootMsg {
	if r == nil {
		return nil
	}
	id := r.getID().String()
	switch root := r.(type) {
	case One:
		return &RootMsg{K: OneK, ID: id}
	case TpRef:
		return &RootMsg{K: RefK, ID: id}
	case Tensor:
		return &RootMsg{
			K:     TensorK,
			ID:    id,
			Value: MsgFromRoot(root.S),
			State: MsgFromRoot(root.T),
		}
	case Lolli:
		return &RootMsg{
			K:     LolliK,
			ID:    id,
			Value: MsgFromRoot(root.S),
			State: MsgFromRoot(root.T),
		}
	case With:
		choices := make([]ChoiceMsg, len(root.Choices))
		for i, l := range maps.Keys(root.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), State: MsgFromRoot(root.Choices[l])}
		}
		return &RootMsg{K: WithK, ID: id, Choices: choices}
	case Plus:
		choices := make([]ChoiceMsg, len(root.Choices))
		for i, l := range maps.Keys(root.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), State: MsgFromRoot(root.Choices[l])}
		}
		return &RootMsg{K: PlusK, ID: id, Choices: choices}
	default:
		panic(ErrUnexpectedRoot(r))
	}
}

func MsgToRoot(mto *RootMsg) (Root, error) {
	if mto == nil {
		return nil, nil
	}
	id, err := id.String[ID](mto.ID)
	if err != nil {
		return nil, err
	}
	switch mto.K {
	case OneK:
		return One{ID: id}, nil
	case RefK:
		return TpRef{ID: id, Name: mto.Name}, nil
	case TensorK:
		v, err := MsgToRoot(mto.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToRoot(mto.State)
		if err != nil {
			return nil, err
		}
		return Tensor{ID: id, S: v, T: s}, nil
	case LolliK:
		v, err := MsgToRoot(mto.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToRoot(mto.State)
		if err != nil {
			return nil, err
		}
		return Lolli{ID: id, S: v, T: s}, nil
	case WithK:
		choices := make(map[Label]Root, len(mto.Choices))
		for _, mto := range mto.Choices {
			choice, err := MsgToRoot(mto.State)
			if err != nil {
				return nil, err
			}
			choices[Label(mto.Label)] = choice
		}
		return With{ID: id, Choices: choices}, nil
	case PlusK:
		choices := make(map[Label]Root, len(mto.Choices))
		for _, mto := range mto.Choices {
			choice, err := MsgToRoot(mto.State)
			if err != nil {
				return nil, err
			}
			choices[Label(mto.Label)] = choice
		}
		return Plus{ID: id, Choices: choices}, nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func ErrUnexpectedKind(k Kind) error {
	return fmt.Errorf("unextected kind %q", k)
}
