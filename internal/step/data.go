package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"
)

type refData struct {
	K  stepKind `db:"kind" json:"kind,omitempty"`
	ID string   `db:"id" json:"id,omitempty"`
}

type rootData struct {
	ID    string    `db:"id"`
	K     stepKind  `db:"kind"`
	ViaID string    `db:"via_id"`
	Term  *termData `db:"payload"`
}

type stepKind int

const (
	proc = iota
	msg
	srv
)

type termData struct {
	K     termKind   `json:"kind"`
	Close *closeData `json:"close,omitempty"`
	Wait  *waitData  `json:"wait,omitempty"`
	Send  *sendData  `json:"send,omitempty"`
	Recv  *recvData  `json:"recv,omitempty"`
	Lab   *labData   `json:"lab,omitempty"`
	Case  *caseData  `json:"case,omitempty"`
}

type closeData struct {
	A string `json:"a"`
}

type waitData struct {
	X    string    `json:"x"`
	Cont *termData `json:"cont"`
}

type sendData struct {
	A string `json:"a"`
	B string `json:"b"`
}

type recvData struct {
	X    string    `json:"x"`
	Y    string    `json:"y"`
	Cont *termData `json:"cont"`
}

type labData struct {
	C     string `json:"c"`
	Label string `json:"label"`
}

type caseData struct {
	Z     string               `json:"z"`
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
		return &rootData{
			K:    proc,
			ID:   root.ID.String(),
			Term: dataFromTerm(root.Term),
		}, nil
	case MsgRoot:
		return &rootData{
			K:     msg,
			ID:    root.ID.String(),
			ViaID: root.ViaID.String(),
			Term:  dataFromValue(root.Val),
		}, nil
	case SrvRoot:
		return &rootData{
			K:     srv,
			ID:    root.ID.String(),
			ViaID: root.ViaID.String(),
			Term:  dataFromCont(root.Cont),
		}, nil
	default:
		panic(ErrUnexpectedStep(root))
	}
}

func dataToRoot(dto *rootData) (root, error) {
	if dto == nil {
		return nil, nil
	}
	rootID, err := id.StringTo(dto.ID)
	if err != nil {
		return nil, err
	}
	viaID, err := id.StringTo(dto.ViaID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case proc:
		term, err := dataToTerm(dto.Term)
		if err != nil {
			return nil, err
		}
		return ProcRoot{ID: rootID, Term: term}, nil
	case msg:
		val, err := dataToValue(dto.Term)
		if err != nil {
			return nil, err
		}
		return MsgRoot{ID: rootID, ViaID: viaID, Val: val}, nil
	case srv:
		cont, err := dataToCont(dto.Term)
		if err != nil {
			return nil, err
		}
		return SrvRoot{ID: rootID, ViaID: viaID, Cont: cont}, nil
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
			Close: &closeData{id.StringFrom(term.A)},
		}
	case WaitSpec:
		return &termData{
			K:    wait,
			Wait: &waitData{id.StringFrom(term.X), dataFromTerm(term.Cont)},
		}
	case SendSpec:
		return &termData{
			K:    send,
			Send: &sendData{id.StringFrom(term.A), id.StringFrom(term.B)},
		}
	case RecvSpec:
		return &termData{
			K: recv,
			Recv: &recvData{
				id.StringFrom(term.X),
				id.StringFrom(term.Y),
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
		a, err := id.StringTo(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case wait:
		x, err := id.StringTo(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case send:
		a, err := id.StringTo(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := id.StringTo(dto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case recv:
		x, err := id.StringTo(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := id.StringTo(dto.Recv.Y)
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
			Close: &closeData{id.StringFrom(val.A)},
		}
	case SendSpec:
		return &termData{
			K:    send,
			Send: &sendData{id.StringFrom(val.A), id.StringFrom(val.B)},
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
		a, err := id.StringTo(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case send:
		a, err := id.StringTo(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := id.StringTo(dto.Send.B)
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
			Wait: &waitData{id.StringFrom(cont.X), dataFromTerm(cont.Cont)},
		}
	case RecvSpec:
		return &termData{
			K:    recv,
			Recv: &recvData{id.StringFrom(cont.X), id.StringFrom(cont.Y), dataFromTerm(cont.Cont)},
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
		x, err := id.StringTo(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case recv:
		x, err := id.StringTo(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := id.StringTo(dto.Recv.Y)
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
