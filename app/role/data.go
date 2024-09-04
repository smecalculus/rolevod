package role

import (
	"errors"

	"smecalculus/rolevod/lib/core"
)

type roleRootData struct {
	ID       string                  `db:"id"`
	Name     string                  `db:"name"`
	Children []roleTeaserData        `db:"-"`
	States   map[string]state        `db:"-"`
	Trs      map[string][]transition `db:"-"`
}

type roleTeaserData struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type kinshipRootData struct {
	Parent   roleTeaserData
	Children []roleTeaserData
}

type kind int

const (
	oneK = iota
	refK
	tensorK
	lolliK
	withK
	plusK
)

type state struct {
	K    kind   `db:"kind"`
	ID   string `db:"id"`
	Name string `db:"name"`
}

type transition struct {
	FromID string `db:"from_id"`
	ToID   string `db:"to_id"`
	MsgID  string `db:"msg_id"`
	MsgKey string `db:"msg_key"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRoleTeaser    func(roleTeaserData) (RoleTeaser, error)
	DataFromRoleTeaser  func(RoleTeaser) roleTeaserData
	DataToRoleTeasers   func([]roleTeaserData) ([]RoleTeaser, error)
	DataFromRoleTeasers func([]RoleTeaser) []roleTeaserData
	DataToRoleRoots     func([]roleRootData) ([]RoleRoot, error)
	DataFromRoleRoots   func([]RoleRoot) []roleRootData
	// kinship
	DataToKinshipRoot   func(kinshipRootData) (KinshipRoot, error)
	DataFromKinshipRoot func(KinshipRoot) kinshipRootData
)

func dataToRoleRoot(dto roleRootData) (RoleRoot, error) {
	id, err := core.FromString[Role](dto.ID)
	if err != nil {
		return RoleRoot{}, nil
	}
	children, err := DataToRoleTeasers(dto.Children)
	if err != nil {
		return RoleRoot{}, nil
	}
	return RoleRoot{
		ID:       id,
		Name:     dto.Name,
		Children: children,
		St:       dataToStype(dto, dto.States[dto.ID]),
	}, nil
}

func dataFromRoleRoot(root RoleRoot) roleRootData {
	data := &roleRootData{
		ID:       root.ID.String(),
		Name:     root.Name,
		Children: DataFromRoleTeasers(root.Children),
		States:   map[string]state{},
		Trs:      map[string][]transition{},
	}
	if root.St == nil {
		return *data
	}
	state := dataFromStype(data, root.St)
	data.States[state.ID] = state
	return *data
}

func dataFromStype(data *roleRootData, stype Stype) state {
	switch st := stype.(type) {
	case One:
		return state{K: oneK, ID: st.ID.String()}
	case TpRef:
		return state{K: refK, ID: st.ID.String(), Name: st.Name}
	case Tensor:
		from := state{K: tensorK, ID: st.ID.String()}
		msg := dataFromStype(data, st.S)
		to := dataFromStype(data, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		data.States[to.ID] = to
		data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		return from
	case Lolli:
		from := state{K: lolliK, ID: st.ID.String()}
		msg := dataFromStype(data, st.S)
		to := dataFromStype(data, st.T)
		tr := transition{FromID: from.ID, ToID: to.ID, MsgID: msg.ID}
		data.States[to.ID] = to
		data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		return from
	case With:
		from := state{K: withK, ID: st.ID.String()}
		for l, st := range st.Chs {
			to := dataFromStype(data, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			data.States[to.ID] = to
			data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		}
		return from
	case Plus:
		from := state{K: plusK, ID: st.ID.String()}
		for l, st := range st.Chs {
			to := dataFromStype(data, st)
			tr := transition{FromID: from.ID, ToID: to.ID, MsgKey: string(l)}
			data.States[to.ID] = to
			data.Trs[from.ID] = append(data.Trs[from.ID], tr)
		}
		return from
	default:
		panic(ErrUnexpectedSt)
	}
}

func dataToStype(data roleRootData, from state) Stype {
	id, err := core.FromString[Role](from.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch from.K {
	case oneK:
		return One{ID: id}
	case refK:
		return TpRef{ID: id, Name: from.Name}
	case tensorK:
		tr := data.Trs[from.ID][0]
		return Tensor{
			ID: id,
			S:  dataToStype(data, data.States[tr.MsgID]),
			T:  dataToStype(data, data.States[tr.ToID]),
		}
	case lolliK:
		tr := data.Trs[from.ID][0]
		return Lolli{
			ID: id,
			S:  dataToStype(data, data.States[tr.MsgID]),
			T:  dataToStype(data, data.States[tr.ToID]),
		}
	case withK:
		st := With{ID: id}
		for _, tr := range data.Trs[from.ID] {
			st.Chs[Label(tr.MsgKey)] = dataToStype(data, data.States[tr.ToID])
		}
		return st
	case plusK:
		st := Plus{ID: id}
		for _, tr := range data.Trs[from.ID] {
			st.Chs[Label(tr.MsgKey)] = dataToStype(data, data.States[tr.ToID])
		}
		return st
	default:
		panic(ErrUnexpectedSt)
	}
}

var (
	errInvalidID = errors.New("invalid state id")
)
