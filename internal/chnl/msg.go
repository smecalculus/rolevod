package chnl

import (
	valid "github.com/go-ozzo/ozzo-validation/v4"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type SpecMsg struct {
	Name string        `json:"name"`
	St   *state.RefMsg `json:"state"`
}

func (mto SpecMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.Name, valid.Required, valid.Length(1, 64)),
		valid.Field(&mto.St, valid.Required),
	)
}

type RefMsg struct {
	ID   string `json:"id" param:"id"`
	Name string `json:"name"`
}

func (mto RefMsg) Validate() error {
	return valid.ValidateStruct(&mto,
		valid.Field(&mto.ID, valid.Required, valid.Length(20, 20)),
		valid.Field(&mto.Name, valid.Required, valid.Length(1, 64)),
	)
}

type RootMsg struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	St   *state.RefMsg `json:"state"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:String.*
// goverter:extend smecalculus/rolevod/lib/ak:String.*
// goverter:extend smecalculus/rolevod/internal/state:Msg.*
var (
	MsgToSpec    func(SpecMsg) (Spec, error)
	MsgFromSpec  func(Spec) SpecMsg
	MsgToRef     func(RefMsg) (Ref, error)
	MsgFromRef   func(Ref) RefMsg
	MsgFromRoot  func(Root) RootMsg
	MsgFromRoots func([]Root) []RootMsg
)

func MsgFromRefMap(refs map[Sym]ID) []RefMsg {
	var mtos []RefMsg
	for name, rid := range refs {
		mtos = append(mtos, RefMsg{rid.String(), string(name)})
	}
	return mtos
}

func MsgToRefMap(mtos []RefMsg) (map[Sym]ID, error) {
	refs := make(map[Sym]ID, len(mtos))
	for _, mto := range mtos {
		rid, err := id.StringToID(mto.ID)
		if err != nil {
			return nil, err
		}
		refs[Sym(mto.Name)] = rid
	}
	return refs, nil
}
