package proc

import (
	"smecalculus/rolevod/internal/step"
)

type modData struct {
	PoolID string
	Bnds   []bndData
	Steps  []step.RootData
	Rev    int
}

type bndData struct {
	ProcID  string
	ProcPH  string
	ChnlID  string
	StateID string
	Rev     int
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
var (
	DataFromMod func(Mod) modData
	DataFromBnd func(Binding) bndData
)
