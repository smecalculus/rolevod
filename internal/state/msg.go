package state

import (
	"fmt"

	"golang.org/x/exp/maps"

	valid "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"
)

type SpecMsg struct {
	K      Kind      `json:"kind"`
	Recur  *RecurMsg `json:"recur,omitempty"`
	Tensor *ProdMsg  `json:"tensor,omitempty"`
	Lolli  *ProdMsg  `json:"lolli,omitempty"`
	With   *SumMsg   `json:"with,omitempty"`
	Plus   *SumMsg   `json:"plus,omitempty"`
}

func (mto SpecMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.K, valid.Required, valid.In(One, Recur, Tensor, Lolli, With, Plus)),
		valid.Field(&mto.Recur, valid.Required.When(mto.K == Recur)),
		valid.Field(&mto.Tensor, valid.Required.When(mto.K == Tensor)),
		valid.Field(&mto.Lolli, valid.Required.When(mto.K == Lolli)),
		valid.Field(&mto.With, valid.Required.When(mto.K == With)),
		valid.Field(&mto.Plus, valid.Required.When(mto.K == Plus)),
	)
}

type RecurMsg struct {
	Name string `json:"name"`
	ToID string `json:"to_id"`
}

type ProdMsg struct {
	Value *SpecMsg `json:"value"`
	State *SpecMsg `json:"state"`
}

func (mto ProdMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.Value, valid.Required),
		valid.Field(&mto.State, valid.Required),
	)
}

type SumMsg struct {
	Choices []ChoiceMsg `json:"choices"`
}

func (mto SumMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.Choices, valid.Required, valid.Length(1, 20)),
	)
}

type ChoiceMsg struct {
	Label string   `json:"label"`
	State *SpecMsg `json:"state"`
}

func (mto ChoiceMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.Label, valid.Required, valid.Length(1, 36)),
		valid.Field(&mto.State, valid.Required),
	)
}

type RefMsg struct {
	ID string `param:"id" json:"id"`
	K  Kind   `json:"kind"`
}

func (mto RefMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.ID, valid.Required, valid.Length(20, 20)),
		valid.Field(&mto.K, valid.Required, valid.In(One, Recur, Tensor, Lolli, With, Plus)),
	)
}

type Kind string

const (
	One    = Kind("one")
	Recur  = Kind("recur")
	Tensor = Kind("tensor")
	Lolli  = Kind("lolli")
	With   = Kind("with")
	Plus   = Kind("plus")
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
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
	case RecurSpec:
		return &SpecMsg{
			K:     Recur,
			Recur: &RecurMsg{ToID: spec.ToID.String(), Name: spec.Name}}
	case TensorSpec:
		return &SpecMsg{
			K: Tensor,
			Tensor: &ProdMsg{
				Value: MsgFromSpec(spec.B),
				State: MsgFromSpec(spec.C),
			},
		}
	case LolliSpec:
		return &SpecMsg{
			K: Lolli,
			Lolli: &ProdMsg{
				Value: MsgFromSpec(spec.Y),
				State: MsgFromSpec(spec.Z),
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
	case Recur:
		id, err := id.StringToID(mto.Recur.ToID)
		if err != nil {
			return nil, err
		}
		return RecurSpec{ToID: id, Name: mto.Recur.Name}, nil
	case Tensor:
		v, err := MsgToSpec(mto.Tensor.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Tensor.State)
		if err != nil {
			return nil, err
		}
		return TensorSpec{B: v, C: s}, nil
	case Lolli:
		v, err := MsgToSpec(mto.Lolli.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Lolli.State)
		if err != nil {
			return nil, err
		}
		return LolliSpec{Y: v, Z: s}, nil
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
	id := ref.RID().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return &RefMsg{K: One, ID: id}
	case RecurRef, RecurRoot:
		return &RefMsg{K: Recur, ID: id}
	case TensorRef, TensorRoot:
		return &RefMsg{K: Tensor, ID: id}
	case LolliRef, LolliRoot:
		return &RefMsg{K: Lolli, ID: id}
	case WithRef, WithRoot:
		return &RefMsg{K: With, ID: id}
	case PlusRef, PlusRoot:
		return &RefMsg{K: Plus, ID: id}
	default:
		panic(ErrUnexpectedRef(ref))
	}
}

func MsgToRef(mto *RefMsg) (Ref, error) {
	if mto == nil {
		return nil, nil
	}
	rid, err := id.StringToID(mto.ID)
	if err != nil {
		return nil, err
	}
	switch mto.K {
	case One:
		return OneRef{rid}, nil
	case Recur:
		return RecurRef{rid}, nil
	case Tensor:
		return TensorRef{rid}, nil
	case Lolli:
		return LolliRef{rid}, nil
	case With:
		return WithRef{rid}, nil
	case Plus:
		return PlusRef{rid}, nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func ErrUnexpectedKind(k Kind) error {
	return fmt.Errorf("unextected kind %q", k)
}
