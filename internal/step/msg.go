package step

import (
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
	Step    *RootMsg    `json:"step,omitempty"`
	Choices []ChoiceMsg `json:"choices,omitempty"`
}

type ChoiceMsg struct {
	Label string  `json:"label"`
	Step  RootMsg `json:"step"`
}

type Kind string

const (
	FwdK   = Kind("fwd")
	SpawnK = Kind("spawn")
	CloseK = Kind("close")
	WaitK  = Kind("wait")
	RefK   = Kind("ref")
	SendK  = Kind("send")
	RecvK  = Kind("recv")
	LabK   = Kind("lab")
	CaseK  = Kind("case")
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
	switch step := root.(type) {
	case nil:
		return nil
	case *ExpRef:
		return &RootMsg{K: RefK, ID: step.ID.String()}
	case *Close:
		return &RootMsg{K: CloseK, ID: step.ID.String()}
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
	case CloseK:
		return &Close{ID: id}, nil
	case RefK:
		return &ExpRef{ID: id, Name: mto.Name}, nil
	default:
		panic(ErrUnexpectedSt)
	}
}
