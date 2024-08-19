package dcl

import "smecalculus/rolevod/lib/core"

type rootData struct {
	ID   string
	Name string
	St   stypeData
}

type stypeData struct {
	Sts []stateData
	Trs []transitionData
}

type tag int

const (
	oneT = iota
	refT
	tensorT
)

type stateData struct {
	ID    int
	T     tag
	RefID string
}

type transitionData struct {
	T    tag
	From stateData
	To   stateData
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToTpRoot   func(rootData) (TpRoot, error)
	DataFromTpRoot func(TpRoot) rootData
)

func dataFromStype(stype Stype) stypeData {
	switch st := stype.(type) {
	case One:
		sd := stateData{T: oneT}
		return stypeData{Sts: []stateData{sd}}
	case TpRef:
		sd := stateData{T: refT, RefID: core.ToString[AR](st.ID)}
		return stypeData{Sts: []stateData{sd}}
	case Tensor:
		sd := stateData{T: tensorT}
		td := transitionData{
			T:    tensorT,
			From: sd,
		}
		return stypeData{
			Sts: []stateData{sd},
			Trs: []transitionData{td},
		}
	default:
		panic(ErrUnexpectedSt)
	}
}
