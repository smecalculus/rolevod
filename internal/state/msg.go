package state

import (
	"fmt"

	"golang.org/x/exp/maps"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
)

type SpecMsg struct {
	K      Kind      `json:"kind"`
	Recur  *RecurMsg `json:"recur,omitempty"`
	Tensor *ProdMsg  `json:"tensor,omitempty"`
	Lolli  *ProdMsg  `json:"lolli,omitempty"`
	Plus   *SumMsg   `json:"plus,omitempty"`
	With   *SumMsg   `json:"with,omitempty"`
}

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.K, kindRequired...),
		validation.Field(&mto.Recur, validation.Required.When(mto.K == Recur)),
		validation.Field(&mto.Tensor, validation.Required.When(mto.K == Tensor)),
		validation.Field(&mto.Lolli, validation.Required.When(mto.K == Lolli)),
		validation.Field(&mto.Plus, validation.Required.When(mto.K == Plus)),
		validation.Field(&mto.With, validation.Required.When(mto.K == With)),
	)
}

type RecurMsg struct {
	Name string `json:"name"`
	ToID string `json:"to_id"`
}

type ProdMsg struct {
	Value SpecMsg `json:"value"`
	Cont  SpecMsg `json:"cont"`
}

func (mto ProdMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Value, validation.Required),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type SumMsg struct {
	Choices []ChoiceMsg `json:"choices"`
}

func (mto SumMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Choices,
			validation.Required,
			validation.Length(1, 10),
			validation.Each(validation.Required),
		),
	)
}

type ChoiceMsg struct {
	Label string  `json:"label"`
	Cont  SpecMsg `json:"cont"`
}

func (mto ChoiceMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.Label, core.NameRequired...),
		validation.Field(&mto.Cont, validation.Required),
	)
}

type RefMsg struct {
	ID string `json:"id" param:"id"`
	K  Kind   `json:"kind"`
}

func (mto RefMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.ID, id.Required...),
		validation.Field(&mto.K, kindRequired...),
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

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(One, Recur, Tensor, Lolli, Plus, With),
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend Msg.*
var (
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromSpec(s Spec) SpecMsg {
	switch spec := s.(type) {
	case OneSpec:
		return SpecMsg{K: One}
	case MenSpec:
		return SpecMsg{
			K:     Recur,
			Recur: &RecurMsg{ToID: spec.StID.String(), Name: spec.Name}}
	case TensorSpec:
		return SpecMsg{
			K: Tensor,
			Tensor: &ProdMsg{
				Value: MsgFromSpec(spec.B),
				Cont:  MsgFromSpec(spec.C),
			},
		}
	case LolliSpec:
		return SpecMsg{
			K: Lolli,
			Lolli: &ProdMsg{
				Value: MsgFromSpec(spec.Y),
				Cont:  MsgFromSpec(spec.Z),
			},
		}
	case WithSpec:
		choices := make([]ChoiceMsg, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), Cont: MsgFromSpec(spec.Choices[l])}
		}
		return SpecMsg{K: With, With: &SumMsg{Choices: choices}}
	case PlusSpec:
		choices := make([]ChoiceMsg, len(spec.Choices))
		for i, l := range maps.Keys(spec.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), Cont: MsgFromSpec(spec.Choices[l])}
		}
		return SpecMsg{K: Plus, Plus: &SumMsg{Choices: choices}}
	default:
		panic(ErrUnexpectedSpec(s))
	}
}

func MsgToSpec(mto SpecMsg) (Spec, error) {
	switch mto.K {
	case One:
		return OneSpec{}, nil
	case Recur:
		id, err := id.StringToID(mto.Recur.ToID)
		if err != nil {
			return nil, err
		}
		return MenSpec{StID: id, Name: mto.Recur.Name}, nil
	case Tensor:
		v, err := MsgToSpec(mto.Tensor.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Tensor.Cont)
		if err != nil {
			return nil, err
		}
		return TensorSpec{B: v, C: s}, nil
	case Lolli:
		v, err := MsgToSpec(mto.Lolli.Value)
		if err != nil {
			return nil, err
		}
		s, err := MsgToSpec(mto.Lolli.Cont)
		if err != nil {
			return nil, err
		}
		return LolliSpec{Y: v, Z: s}, nil
	case With:
		choices := make(map[Label]Spec, len(mto.With.Choices))
		for _, ch := range mto.With.Choices {
			choice, err := MsgToSpec(ch.Cont)
			if err != nil {
				return nil, err
			}
			choices[Label(ch.Label)] = choice
		}
		return WithSpec{Choices: choices}, nil
	case Plus:
		choices := make(map[Label]Spec, len(mto.Plus.Choices))
		for _, ch := range mto.Plus.Choices {
			choice, err := MsgToSpec(ch.Cont)
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

func MsgFromRef(ref Ref) RefMsg {
	id := ref.RID().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return RefMsg{K: One, ID: id}
	case MenRef, MenRoot:
		return RefMsg{K: Recur, ID: id}
	case TensorRef, TensorRoot:
		return RefMsg{K: Tensor, ID: id}
	case LolliRef, LolliRoot:
		return RefMsg{K: Lolli, ID: id}
	case PlusRef, PlusRoot:
		return RefMsg{K: Plus, ID: id}
	case WithRef, WithRoot:
		return RefMsg{K: With, ID: id}
	default:
		panic(ErrUnexpectedRef(ref))
	}
}

func MsgToRef(mto RefMsg) (Ref, error) {
	rid, err := id.StringToID(mto.ID)
	if err != nil {
		return nil, err
	}
	switch mto.K {
	case One:
		return OneRef{rid}, nil
	case Recur:
		return MenRef{rid}, nil
	case Tensor:
		return TensorRef{rid}, nil
	case Lolli:
		return LolliRef{rid}, nil
	case Plus:
		return PlusRef{rid}, nil
	case With:
		return WithRef{rid}, nil
	default:
		panic(ErrUnexpectedKind(mto.K))
	}
}

func ErrUnexpectedKind(k Kind) error {
	return fmt.Errorf("unextected kind %q", k)
}
