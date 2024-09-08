package step

import (
	"errors"

	"smecalculus/rolevod/lib/id"
)

type rootData struct {
	ID    string
	Steps map[string]step
	Trs   map[string][]transition
}

type kind int

const (
	fwdK = iota
	spawnK
	closeK
	waitK
	refK
	sendK
	recvK
	labK
	caseK
)

type step struct {
	K    kind   `db:"kind"`
	ID   string `db:"id"`
	Name string `db:"name"`
}

type transition struct {
	FromID string `db:"from_id"`
	ToID   string `db:"to_id"`
	ValID  string `db:"val_id"`
	MsgKey string `db:"msg_key"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRef func(step) (Ref, error)
	// goverter:ignore K Name
	DataFromRef   func(Ref) step
	DataToRefs    func([]step) ([]Ref, error)
	DataFromRefs  func([]Ref) []step
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) []*rootData
)

func dataToRoot(dto rootData) Root {
	return dataToStep(dto, dto.Steps[dto.ID])
}

func dataFromRoot(root Root) *rootData {
	if root == nil {
		return nil
	}
	dto := &rootData{
		ID:    root.get().String(),
		Steps: map[string]step{},
		Trs:   map[string][]transition{},
	}
	step := dataFromStep(dto, root)
	dto.Steps[step.ID] = step
	return dto
}

func dataFromStep(dto *rootData, root Root) step {
	switch st := root.(type) {
	case *ExpRef:
		return step{K: refK, ID: st.ID.String()}
	case *Close:
		return step{K: closeK, ID: st.ID.String()}
	case *Wait:
		return step{K: waitK, ID: st.ID.String()}
	default:
		panic(ErrUnexpectedSt)
	}
}

func dataToStep(root rootData, step step) Root {
	id, err := id.String[ID](step.ID)
	if err != nil {
		panic(errInvalidID)
	}
	switch step.K {
	case refK:
		return &ExpRef{ID: id}
	case closeK:
		return &Close{ID: id}
	default:
		panic(ErrUnexpectedSt)
	}
}

var (
	errInvalidID = errors.New("invalid step id")
)
