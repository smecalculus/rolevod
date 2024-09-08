package role

import (
	"golang.org/x/exp/maps"

	"smecalculus/rolevod/lib/id"
)

type RoleSpecMsg struct {
	Name string    `json:"name"`
	St   *StypeMsg `json:"st,omitempty"`
}

type RefMsg struct {
	ID string `param:"id" query:"id" json:"id"`
}

type RoleRootMsg struct {
	ID       string       `param:"id" json:"id"`
	Name     string       `json:"name"`
	Children []RoleRefMsg `json:"children"`
	St       *StypeMsg    `json:"st,omitempty"`
}

type RoleRefMsg struct {
	ID   string `param:"id" json:"id"`
	Name string `query:"name" json:"name"`
}

type KinshipSpecMsg struct {
	Parent   string   `param:"id" json:"parent"`
	Children []string `json:"children"`
}

type KinshipRootMsg struct {
	Parent   RoleRefMsg   `json:"parent"`
	Children []RoleRefMsg `json:"children"`
}

type StypeMsg struct {
	K    Kind        `json:"kind"`
	ID   string      `json:"id"`
	Name string      `json:"name,omitempty"`
	M    *StypeMsg   `json:"message,omitempty"`
	S    *StypeMsg   `json:"session,omitempty"`
	Chs  []ChoiceMsg `json:"choices,omitempty"`
}

type ChoiceMsg struct {
	L string   `json:"label"`
	S StypeMsg `json:"session"`
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

type OneMsg struct {
	K  Kind   `json:"kind"`
	ID string `json:"id"`
}

type TpRefMsg struct {
	K    Kind   `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend msg.*
var (
	// role
	MsgToRoleSpec    func(RoleSpecMsg) (RoleSpec, error)
	MsgFromRoleSpec  func(RoleSpec) RoleSpecMsg
	MsgFromRoleRoot  func(RoleRoot) RoleRootMsg
	MsgToRoleRoot    func(RoleRootMsg) (RoleRoot, error)
	MsgFromRoleRoots func([]RoleRoot) []RoleRootMsg
	MsgToRoleRoots   func([]RoleRootMsg) ([]RoleRoot, error)
	MsgFromRoleRef   func(RoleRef) RoleRefMsg
	MsgToRoleRef     func(RoleRefMsg) (RoleRef, error)
	MsgFromRoleRefs  func([]RoleRef) []RoleRefMsg
	MsgToRoleRefs    func([]RoleRefMsg) ([]RoleRef, error)
	// kinship
	MsgFromKinshipSpec func(KinshipSpec) KinshipSpecMsg
	MsgToKinshipSpec   func(KinshipSpecMsg) (KinshipSpec, error)
	MsgFromKinshipRoot func(KinshipRoot) KinshipRootMsg
	MsgToKinshipRoot   func(KinshipRootMsg) (KinshipRoot, error)
)

func msgFromStype(stype Stype) *StypeMsg {
	switch st := stype.(type) {
	case nil:
		return nil
	case One:
		return &StypeMsg{K: OneK, ID: st.ID.String()}
	case TpRef:
		return &StypeMsg{K: RefK, ID: st.ID.String(), Name: st.Name}
	case Tensor:
		return &StypeMsg{
			K:  TensorK,
			ID: st.ID.String(),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case Lolli:
		return &StypeMsg{
			K:  LolliK,
			ID: st.ID.String(),
			M:  msgFromStype(st.S),
			S:  msgFromStype(st.T),
		}
	case With:
		sts := make([]ChoiceMsg, len(st.Choices))
		for i, l := range maps.Keys(st.Choices) {
			sts[i] = ChoiceMsg{L: string(l), S: *msgFromStype(st.Choices[l])}
		}
		return &StypeMsg{K: WithK, ID: st.ID.String(), Chs: sts}
	case Plus:
		sts := make([]ChoiceMsg, len(st.Choices))
		for i, l := range maps.Keys(st.Choices) {
			sts[i] = ChoiceMsg{L: string(l), S: *msgFromStype(st.Choices[l])}
		}
		return &StypeMsg{K: PlusK, ID: st.ID.String(), Chs: sts}
	default:
		panic(ErrUnexpectedSt)
	}
}

func msgToStype(msg *StypeMsg) (Stype, error) {
	if msg == nil {
		return nil, nil
	}
	id, err := id.String[ID](msg.ID)
	if err != nil {
		return nil, err
	}
	switch msg.K {
	case OneK:
		return One{ID: id}, nil
	case RefK:
		return TpRef{ID: id, Name: msg.Name}, nil
	case TensorK:
		m, err := msgToStype(msg.M)
		if err != nil {
			return nil, err
		}
		s, err := msgToStype(msg.S)
		if err != nil {
			return nil, err
		}
		return Tensor{ID: id, S: m, T: s}, nil
	case LolliK:
		m, err := msgToStype(msg.M)
		if err != nil {
			return nil, err
		}
		s, err := msgToStype(msg.S)
		if err != nil {
			return nil, err
		}
		return Lolli{ID: id, S: m, T: s}, nil
	case WithK:
		sts := make(map[Label]Stype, len(msg.Chs))
		for _, ch := range msg.Chs {
			st, err := msgToStype(&ch.S)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.L)] = st
		}
		return With{ID: id, Choices: sts}, nil
	case PlusK:
		sts := make(map[Label]Stype, len(msg.Chs))
		for _, ch := range msg.Chs {
			st, err := msgToStype(&ch.S)
			if err != nil {
				return nil, err
			}
			sts[Label(ch.L)] = st
		}
		return Plus{ID: id, Choices: sts}, nil
	default:
		panic(ErrUnexpectedSt)
	}
}
