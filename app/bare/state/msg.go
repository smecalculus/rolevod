package state

import (
	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/id"
)

type RefMsg struct {
	ID string `param:"id" json:"id"`
}

type RootMsg struct {
	K       Kind        `json:"kind"`
	ID      string      `json:"id"`
	Name    string      `json:"name,omitempty"`
	Value   *RootMsg    `json:"value,omitempty"`
	State   *RootMsg    `json:"state,omitempty"`
	Choices []ChoiceMsg `json:"choices,omitempty"`
}

type ChoiceMsg struct {
	Label  string  `json:"label"`
	String RootMsg `json:"state"`
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
var (
	MsgFromRef  func(Ref) RefMsg
	MsgToRef    func(RefMsg) (Ref, error)
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromRoot(root Root) *RootMsg {
	switch state := root.(type) {
	case nil:
		return nil
	case *TpRef:
		return &RootMsg{K: RefK, ID: state.ID.String()}
	case *One:
		return &RootMsg{K: OneK, ID: state.ID.String()}
	case *Tensor:
		return &RootMsg{
			K:     TensorK,
			ID:    state.ID.String(),
			Value: MsgFromRoot(state.S),
			State: MsgFromRoot(state.T),
		}
	case *Lolli:
		return &RootMsg{
			K:     LolliK,
			ID:    state.ID.String(),
			Value: MsgFromRoot(state.S),
			State: MsgFromRoot(state.T),
		}
	case *With:
		sts := make([]ChoiceMsg, len(state.Choices))
		for i, l := range maps.Keys(state.Choices) {
			sts[i] = ChoiceMsg{Label: string(l), String: *MsgFromRoot(state.Choices[l])}
		}
		return &RootMsg{K: WithK, ID: state.ID.String(), Choices: sts}
	case *Plus:
		sts := make([]ChoiceMsg, len(state.Choices))
		for i, l := range maps.Keys(state.Choices) {
			sts[i] = ChoiceMsg{Label: string(l), String: *MsgFromRoot(state.Choices[l])}
		}
		return &RootMsg{K: PlusK, ID: state.ID.String(), Choices: sts}
	default:
		panic(ErrUnexpectedSt)
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
		return &One{ID: id}, nil
	case RefK:
		return &TpRef{ID: id, Name: mto.Name}, nil
	case TensorK:
		m, err := MsgToRoot(mto.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToRoot(mto.State)
		if err != nil {
			return nil, err
		}
		return &Tensor{ID: id, S: m, T: s}, nil
	case LolliK:
		m, err := MsgToRoot(mto.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToRoot(mto.State)
		if err != nil {
			return nil, err
		}
		return &Lolli{ID: id, S: m, T: s}, nil
	case WithK:
		sts := make(map[Label]Root, len(mto.Choices))
		for _, ch := range mto.Choices {
			st, err := MsgToRoot(&ch.String)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.Label)] = st
		}
		return &With{ID: id, Choices: sts}, nil
	case PlusK:
		sts := make(map[Label]Root, len(mto.Choices))
		for _, ch := range mto.Choices {
			st, err := MsgToRoot(&ch.String)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.Label)] = st
		}
		return &Plus{ID: id, Choices: sts}, nil
	default:
		panic(ErrUnexpectedSt)
	}
}
