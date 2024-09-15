package step

import (
	"encoding/json"
	"errors"
	"fmt"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"
)

type refData struct {
	K  stepKind `db:"kind" json:"kind,omitempty"`
	ID string   `db:"id" json:"id,omitempty"`
}

type rootData struct {
	K       stepKind `db:"kind"`
	ID      string   `db:"id"`
	ViaID   string   `db:"via_id"`
	PreID   string   `db:"pre_id"`
	Payload string   `db:"payload"`
}

type payload struct {
	K     termKind      `json:"kind"`
	Xa    *chnl.RefData `json:"xa,omitempty"`
	Yb    *chnl.RefData `json:"yb,omitempty"`
	Zc    *chnl.RefData `json:"zc,omitempty"`
	Cont  *payload      `json:"cont,omitempty"`
	Label string        `json:"label,omitempty"`
	Conts []payload     `json:"conts,omitempty"`
	Name  string        `json:"name,omitempty"`
	Ctx   []string      `json:"ctx,omitempty"`
}

type stepKind int

const (
	procK = iota
	msgK
	srvK
)

type termKind int

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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend to.*
// goverter:extend data.*
var (
	DataToRef func(refData) (Ref, error)
	// goverter:ignore K
	DataFromRef    func(Ref) refData
	DataToRefs     func([]refData) ([]Ref, error)
	DataFromRefs   func([]Ref) []refData
	DataToTerms    func([]*payload) ([]Term, error)
	DataFromTerms  func([]Term) []*payload
	DataToValues   func([]*payload) ([]Value, error)
	DataFromValues func([]Value) []*payload
	DataToConts    func([]*payload) ([]Continuation, error)
	DataFromConts  func([]Continuation) []*payload
)

func dataFromRoot(r root) (*rootData, error) {
	switch root := r.(type) {
	case Process:
		pl, err := json.Marshal(dataFromTerm(root.Term))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:       procK,
			ID:      root.ID.String(),
			PreID:   root.PreID.String(),
			Payload: string(pl),
		}, nil
	case Message:
		pl, err := json.Marshal(dataFromValue(root.Val))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:       msgK,
			ID:      root.ID.String(),
			PreID:   root.PreID.String(),
			ViaID:   root.ViaID.String(),
			Payload: string(pl),
		}, nil
	case Service:
		pl, err := json.Marshal(dataFromCont(root.Cont))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:       srvK,
			ID:      root.ID.String(),
			PreID:   root.PreID.String(),
			ViaID:   root.ViaID.String(),
			Payload: string(pl),
		}, nil
	default:
		panic(ErrUnexpectedStep(root))
	}
}

func dataToRoot(dto *rootData) (root, error) {
	rootID, err := id.String[ID](dto.ID)
	if err != nil {
		return nil, err
	}
	preID, err := id.String[ID](dto.PreID)
	if err != nil {
		return nil, err
	}
	viaID, err := id.String[chnl.ID](dto.ViaID)
	if err != nil {
		return nil, err
	}
	var pl payload
	err = json.Unmarshal([]byte(dto.Payload), pl)
	if err != nil {
		return nil, err
	}
	var ref *chnl.RefData
	err = json.Unmarshal([]byte(dto.ViaID), ref)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case procK:
		term, err := dataToTerm(&pl)
		if err != nil {
			return nil, err
		}
		return Process{ID: rootID, PreID: preID, Term: term}, nil
	case msgK:
		val, err := dataToValue(&pl)
		if err != nil {
			return nil, err
		}
		return Message{ID: rootID, PreID: preID, ViaID: viaID, Val: val}, nil
	case srvK:
		cont, err := dataToCont(&pl)
		if err != nil {
			return nil, err
		}
		return Service{ID: rootID, PreID: preID, ViaID: viaID, Cont: cont}, nil
	default:
		panic(ErrUnexpectedStepKind(dto.K))
	}
}

func dataFromTerm(t Term) *payload {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case Close:
		return &payload{
			K:  closeK,
			Xa: chnl.DataFromRef(term.A),
		}
	case Wait:
		return &payload{
			K:    waitK,
			Xa:   chnl.DataFromRef(term.X),
			Cont: dataFromTerm(term.Cont),
		}
	case Send:
		return &payload{
			K:  sendK,
			Xa: chnl.DataFromRef(term.A),
			Yb: chnl.DataFromRef(term.B),
		}
	case Recv:
		return &payload{
			K:  recvK,
			Xa: chnl.DataFromRef(term.X),
			Yb: chnl.DataFromRef(term.Y),
		}
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func dataToTerm(dto *payload) (Term, error) {
	xa, err := chnl.DataToRef(*dto.Xa)
	if err != nil {
		return nil, err
	}
	yb, err := chnl.DataToRef(*dto.Yb)
	if err != nil {
		return nil, err
	}
	cont, err := dataToTerm(dto.Cont)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case closeK:
		return Close{A: xa}, nil
	case waitK:
		return Wait{X: xa, Cont: cont}, nil
	case sendK:
		return Send{A: xa, B: yb}, nil
	case recvK:
		return Recv{X: xa, Y: yb, Cont: cont}, nil
	default:
		panic(ErrUnexpectedTermKind(dto.K))
	}
}

func dataFromValue(v Value) *payload {
	switch val := v.(type) {
	case Close:
		return &payload{K: closeK, Xa: chnl.DataFromRef(val.A)}
	case Send:
		return &payload{K: sendK, Xa: chnl.DataFromRef(val.A), Yb: chnl.DataFromRef(val.B)}
	default:
		panic(ErrUnexpectedValue(val))
	}
}

func dataToValue(dto *payload) (Value, error) {
	xa, err := chnl.DataToRef(*dto.Xa)
	if err != nil {
		return nil, err
	}
	yb, err := chnl.DataToRef(*dto.Yb)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case closeK:
		return Close{A: xa}, nil
	case sendK:
		return Send{A: xa, B: yb}, nil
	default:
		panic(ErrUnexpectedTermKind(dto.K))
	}
}

func dataFromCont(c Continuation) *payload {
	switch cont := c.(type) {
	case Wait:
		return &payload{
			K:    waitK,
			Xa:   chnl.DataFromRef(cont.X),
			Cont: dataFromTerm(cont.Cont),
		}
	case Recv:
		return &payload{
			K:  recvK,
			Xa: chnl.DataFromRef(cont.X),
			Yb: chnl.DataFromRef(cont.Y),
		}
	default:
		panic(ErrUnexpectedCont(cont))
	}
}

func dataToCont(dto *payload) (Continuation, error) {
	xa, err := chnl.DataToRef(*dto.Xa)
	if err != nil {
		return nil, err
	}
	yb, err := chnl.DataToRef(*dto.Yb)
	if err != nil {
		return nil, err
	}
	cont, err := dataToTerm(dto.Cont)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case waitK:
		return Wait{X: xa, Cont: cont}, nil
	case recvK:
		return Recv{X: xa, Y: yb, Cont: cont}, nil
	default:
		panic(ErrUnexpectedTermKind(dto.K))
	}
}

var (
	errInvalidID = errors.New("invalid step id")
)

func ErrUnexpectedTermKind(k termKind) error {
	return fmt.Errorf("unexpected term kind %v", k)
}

func ErrUnexpectedStepKind(k stepKind) error {
	return fmt.Errorf("unexpected step kind %v", k)
}
