package proc

import (
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/ph"
	"smecalculus/rolevod/lib/rev"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/step"
)

// aka Configuration
type Snap struct {
	ProcID id.ADT
	EPs    map[ph.ADT]EP
	Steps  map[chnl.ID]step.Root
}

type EP struct {
	ProcID   id.ADT
	ChnlPH   ph.ADT
	ChnlID   id.ADT
	StateID  id.ADT
	PrvdID   id.ADT
	PrvdRevs []rev.ADT
	ClntID   id.ADT
	ClntRevs []rev.ADT
}

func ProcPH(ep EP) ph.ADT { return ep.ChnlPH }

type Mod struct {
	PoolID id.ADT
	Bnds   []Binding
	Steps  []step.Root
	Rev    rev.ADT
}

type Binding struct {
	ProcID  id.ADT
	ProcPH  ph.ADT
	ChnlID  id.ADT
	StateID id.ADT
	Rev     rev.ADT
}
