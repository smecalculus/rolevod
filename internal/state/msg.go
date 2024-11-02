package state

import (
	"fmt"

	"golang.org/x/exp/maps"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"
)

type SpecMsg struct {
	ID     string   `json:"id,omitempty"`
	K      Kind     `json:"kind"`
	Link   *LinkMsg `json:"link,omitempty"`
	Tensor *ProdMsg `json:"tensor,omitempty"`
	Lolli  *ProdMsg `json:"lolli,omitempty"`
	Plus   *SumMsg  `json:"plus,omitempty"`
	With   *SumMsg  `json:"with,omitempty"`
}

func (mto SpecMsg) Validate() error {
	return validation.ValidateStruct(&mto,
		validation.Field(&mto.K, kindRequired...),
		validation.Field(&mto.Link, validation.Required.When(mto.K == Link)),
		validation.Field(&mto.Tensor, validation.Required.When(mto.K == Tensor)),
		validation.Field(&mto.Lolli, validation.Required.When(mto.K == Lolli)),
		validation.Field(&mto.Plus, validation.Required.When(mto.K == Plus)),
		validation.Field(&mto.With, validation.Required.When(mto.K == With)),
	)
}

type LinkMsg struct {
	FQN string `json:"fqn"`
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
	Link   = Kind("link")
	Tensor = Kind("tensor")
	Lolli  = Kind("lolli")
	Plus   = Kind("plus")
	With   = Kind("with")
)

var kindRequired = []validation.Rule{
	validation.Required,
	validation.In(One, Link, Tensor, Lolli, Plus, With),
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend Msg.*
var (
	MsgFromRefs func([]Ref) []RefMsg
	MsgToRefs   func([]RefMsg) ([]Ref, error)
)

func MsgFromRoot(r Root) SpecMsg {
	switch root := r.(type) {
	case OneRoot:
		return SpecMsg{ID: root.ID.String(), K: One}
	case LinkRoot:
		return SpecMsg{
			ID:   root.ID.String(),
			K:    Link,
			Link: &LinkMsg{FQN: sym.StringFromSym(root.Role)}}
	case TensorRoot:
		return SpecMsg{
			ID: root.ID.String(),
			K:  Tensor,
			Tensor: &ProdMsg{
				Value: MsgFromRoot(root.B),
				Cont:  MsgFromRoot(root.C),
			},
		}
	case LolliRoot:
		return SpecMsg{
			ID: root.ID.String(),
			K:  Lolli,
			Lolli: &ProdMsg{
				Value: MsgFromRoot(root.Y),
				Cont:  MsgFromRoot(root.Z),
			},
		}
	case WithRoot:
		choices := make([]ChoiceMsg, len(root.Choices))
		for i, l := range maps.Keys(root.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), Cont: MsgFromRoot(root.Choices[l])}
		}
		return SpecMsg{
			ID:   root.ID.String(),
			K:    With,
			With: &SumMsg{Choices: choices},
		}
	case PlusRoot:
		choices := make([]ChoiceMsg, len(root.Choices))
		for i, l := range maps.Keys(root.Choices) {
			choices[i] = ChoiceMsg{Label: string(l), Cont: MsgFromRoot(root.Choices[l])}
		}
		return SpecMsg{
			ID:   root.ID.String(),
			K:    Plus,
			Plus: &SumMsg{Choices: choices},
		}
	default:
		panic(ErrRootTypeUnexpected(r))
	}
}

func MsgFromSpec(s Spec) SpecMsg {
	switch spec := s.(type) {
	case OneSpec:
		return SpecMsg{K: One}
	case LinkSpec:
		return SpecMsg{
			K:    Link,
			Link: &LinkMsg{FQN: sym.StringFromSym(spec.Role)}}
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
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToRoot(dto SpecMsg) (Root, error) {
	rid, err := id.StringToID(dto.ID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case One:
		return OneRoot{rid}, nil
	case Link:
		return LinkRoot{ID: rid, Role: sym.StringToSym(dto.Link.FQN)}, nil
	case Tensor:
		val, err := MsgToRoot(dto.Tensor.Value)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToRoot(dto.Tensor.Cont)
		if err != nil {
			return nil, err
		}
		return TensorRoot{ID: rid, B: val, C: cont}, nil
	case Lolli:
		val, err := MsgToRoot(dto.Lolli.Value)
		if err != nil {
			return nil, err
		}
		cont, err := MsgToRoot(dto.Lolli.Cont)
		if err != nil {
			return nil, err
		}
		return LolliRoot{ID: rid, Y: val, Z: cont}, nil
	case Plus:
		choices := make(map[core.Label]Root, len(dto.Plus.Choices))
		for _, ch := range dto.Plus.Choices {
			cont, err := MsgToRoot(ch.Cont)
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Label)] = cont
		}
		return PlusRoot{ID: rid, Choices: choices}, nil
	case With:
		choices := make(map[core.Label]Root, len(dto.With.Choices))
		for _, ch := range dto.With.Choices {
			choice, err := MsgToRoot(ch.Cont)
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Label)] = choice
		}
		return WithRoot{ID: rid, Choices: choices}, nil
	default:
		panic(errKindUnexpected(dto.K))
	}
}

func MsgToSpec(mto SpecMsg) (Spec, error) {
	switch mto.K {
	case One:
		return OneSpec{}, nil
	case Link:
		return LinkSpec{Role: sym.StringToSym(mto.Link.FQN)}, nil
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
	case Plus:
		choices := make(map[core.Label]Spec, len(mto.Plus.Choices))
		for _, ch := range mto.Plus.Choices {
			choice, err := MsgToSpec(ch.Cont)
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Label)] = choice
		}
		return PlusSpec{Choices: choices}, nil
	case With:
		choices := make(map[core.Label]Spec, len(mto.With.Choices))
		for _, ch := range mto.With.Choices {
			choice, err := MsgToSpec(ch.Cont)
			if err != nil {
				return nil, err
			}
			choices[core.Label(ch.Label)] = choice
		}
		return WithSpec{Choices: choices}, nil
	default:
		panic(errKindUnexpected(mto.K))
	}
}

func MsgFromRef(ref Ref) RefMsg {
	ident := ref.Ident().String()
	switch ref.(type) {
	case OneRef, OneRoot:
		return RefMsg{K: One, ID: ident}
	case LinkRef, LinkRoot:
		return RefMsg{K: Link, ID: ident}
	case TensorRef, TensorRoot:
		return RefMsg{K: Tensor, ID: ident}
	case LolliRef, LolliRoot:
		return RefMsg{K: Lolli, ID: ident}
	case PlusRef, PlusRoot:
		return RefMsg{K: Plus, ID: ident}
	case WithRef, WithRoot:
		return RefMsg{K: With, ID: ident}
	default:
		panic(ErrRefTypeUnexpected(ref))
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
	case Link:
		return LinkRef{rid}, nil
	case Tensor:
		return TensorRef{rid}, nil
	case Lolli:
		return LolliRef{rid}, nil
	case Plus:
		return PlusRef{rid}, nil
	case With:
		return WithRef{rid}, nil
	default:
		panic(errKindUnexpected(mto.K))
	}
}
func ErrPolarityUnexpected(got Root) error {
	return fmt.Errorf("root polarity unexpected: %v", got.Pol())
}

func ErrPolarityMismatch(a, b Root) error {
	return fmt.Errorf("root polarity mismatch: %v!=%v", a.Pol(), b.Pol())
}

func errKindUnexpected(got Kind) error {
	return fmt.Errorf("kind unexpected: %v", got)
}
