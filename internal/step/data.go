package step

import (
	"encoding/json"
	"fmt"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/lib/id"
)

type refData struct {
	K  stepKind `db:"kind" json:"kind,omitempty"`
	ID string   `db:"id" json:"id,omitempty"`
}

type rootData struct {
	K     stepKind `db:"kind"`
	ID    string   `db:"id"`
	ViaID string   `db:"via_id"`
	PreID string   `db:"pre_id"`
	Term  string   `db:"payload"`
}

type stepKind int

const (
	proc = iota
	msg
	srv
)

type termData struct {
	K termKind `json:"kind"`
	// Xa    *chnl.RefData `json:"xa,omitempty"`
	// Yb    *chnl.RefData `json:"yb,omitempty"`
	// Zc    *chnl.RefData `json:"zc,omitempty"`
	// Cont  *termData     `json:"cont,omitempty"`
	// Label string        `json:"label,omitempty"`
	// Conts []termData    `json:"conts,omitempty"`
	// Name  string        `json:"name,omitempty"`
	// Ctx   []string      `json:"ctx,omitempty"`
	Close *closeData `json:"close,omitempty"`
	Wait  *waitData  `json:"wait,omitempty"`
	Send  *sendData  `json:"send,omitempty"`
	Recv  *recvData  `json:"recv,omitempty"`
	Lab   *labData   `json:"lab,omitempty"`
	Case  *caseData  `json:"case,omitempty"`
}

type closeData struct {
	A chnl.RefData `json:"a"`
}

type waitData struct {
	X    chnl.RefData `json:"x"`
	Cont *termData    `json:"cont"`
}

type sendData struct {
	A chnl.RefData `json:"a"`
	B chnl.RefData `json:"b"`
}

type recvData struct {
	X    chnl.RefData `json:"x"`
	Y    chnl.RefData `json:"y"`
	Cont *termData    `json:"cont"`
}

type labData struct {
	C     chnl.RefData `json:"c"`
	Label string       `json:"label"`
}

type caseData struct {
	Z     chnl.RefData         `json:"z"`
	Conts map[string]*termData `json:"conts"`
}

type termKind int

const (
	fwd = iota
	spawn
	close
	wait
	rec
	send
	recv
	lab
	caseK
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend data.*
var (
	// DataToRef func(refData) (Ref, error)
	// DataFromRef    func(Ref) refData
	// DataToRefs     func([]refData) ([]Ref, error)
	// DataFromRefs   func([]Ref) []refData
	DataToTerms    func([]*termData) ([]Term, error)
	DataFromTerms  func([]Term) []*termData
	DataToValues   func([]*termData) ([]Value, error)
	DataFromValues func([]Value) []*termData
	DataToConts    func([]*termData) ([]Continuation, error)
	DataFromConts  func([]Continuation) []*termData
)

func dataFromRoot(r root) (*rootData, error) {
	if r == nil {
		return nil, nil
	}
	switch root := r.(type) {
	case ProcRoot:
		pl, err := json.Marshal(dataFromTerm(root.Term))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:    proc,
			ID:   root.ID.String(),
			Term: string(pl),
		}, nil
	case MsgRoot:
		pl, err := json.Marshal(dataFromValue(root.Val))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:     msg,
			ID:    root.ID.String(),
			PreID: root.PreID.String(),
			ViaID: root.ViaID.String(),
			Term:  string(pl),
		}, nil
	case SrvRoot:
		pl, err := json.Marshal(dataFromCont(root.Cont))
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:     srv,
			ID:    root.ID.String(),
			PreID: root.PreID.String(),
			ViaID: root.ViaID.String(),
			Term:  string(pl),
		}, nil
	default:
		panic(ErrUnexpectedStep(root))
	}
}

func dataToRoot(dto *rootData) (root, error) {
	if dto == nil {
		return nil, nil
	}
	rootID, err := id.StringFrom(dto.ID)
	if err != nil {
		return nil, err
	}
	preID, err := id.StringFrom(dto.PreID)
	if err != nil {
		return nil, err
	}
	viaID, err := id.StringFrom(dto.ViaID)
	if err != nil {
		return nil, err
	}
	var pl termData
	err = json.Unmarshal([]byte(dto.Term), &pl)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case proc:
		term, err := dataToTerm(&pl)
		if err != nil {
			return nil, err
		}
		return ProcRoot{ID: rootID, Term: term}, nil
	case msg:
		val, err := dataToValue(&pl)
		if err != nil {
			return nil, err
		}
		return MsgRoot{ID: rootID, PreID: preID, ViaID: viaID, Val: val}, nil
	case srv:
		cont, err := dataToCont(&pl)
		if err != nil {
			return nil, err
		}
		return SrvRoot{ID: rootID, PreID: preID, ViaID: viaID, Cont: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func dataFromTerm(t Term) *termData {
	if t == nil {
		return nil
	}
	switch term := t.(type) {
	case CloseSpec:
		return &termData{
			K:     close,
			Close: &closeData{chnl.DataFromRef(term.A)},
		}
	case WaitSpec:
		return &termData{
			K:    wait,
			Wait: &waitData{chnl.DataFromRef(term.X), dataFromTerm(term.Cont)},
		}
	case SendSpec:
		return &termData{
			K:    send,
			Send: &sendData{chnl.DataFromRef(term.A), chnl.DataFromRef(term.B)},
		}
	case RecvSpec:
		return &termData{
			K: recv,
			Recv: &recvData{
				chnl.DataFromRef(term.X),
				chnl.DataFromRef(term.Y),
				dataFromTerm(term.Cont),
			},
		}
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func dataToTerm(dto *termData) (Term, error) {
	if dto == nil {
		return nil, nil
	}
	switch dto.K {
	case close:
		a, err := chnl.DataToRef(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case wait:
		x, err := chnl.DataToRef(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case send:
		a, err := chnl.DataToRef(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := chnl.DataToRef(dto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case recv:
		x, err := chnl.DataToRef(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := chnl.DataToRef(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromValue(v Value) *termData {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case CloseSpec:
		return &termData{
			K:     close,
			Close: &closeData{chnl.DataFromRef(val.A)},
		}
	case SendSpec:
		return &termData{
			K:    send,
			Send: &sendData{chnl.DataFromRef(val.A), chnl.DataFromRef(val.B)},
		}
	default:
		panic(ErrUnexpectedValue(val))
	}
}

func dataToValue(dto *termData) (Value, error) {
	if dto == nil {
		return nil, nil
	}
	switch dto.K {
	case close:
		a, err := chnl.DataToRef(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case send:
		a, err := chnl.DataToRef(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := chnl.DataToRef(dto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromCont(c Continuation) *termData {
	if c == nil {
		return nil
	}
	switch cont := c.(type) {
	case WaitSpec:
		return &termData{
			K:    wait,
			Wait: &waitData{chnl.DataFromRef(cont.X), dataFromTerm(cont.Cont)},
		}
	case RecvSpec:
		return &termData{
			K:    recv,
			Recv: &recvData{chnl.DataFromRef(cont.X), chnl.DataFromRef(cont.Y), dataFromTerm(cont.Cont)},
		}
	default:
		panic(ErrUnexpectedCont(cont))
	}
}

func dataToCont(dto *termData) (Continuation, error) {
	if dto == nil {
		return nil, nil
	}
	switch dto.K {
	case wait:
		x, err := chnl.DataToRef(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case recv:
		x, err := chnl.DataToRef(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := chnl.DataToRef(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func errUnexpectedTermKind(k termKind) error {
	return fmt.Errorf("unexpected term kind %v", k)
}

func errUnexpectedStepKind(k stepKind) error {
	return fmt.Errorf("unexpected step kind %v", k)
}
