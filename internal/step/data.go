package step

import (
	"fmt"

	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/state"
)

type refData struct {
	K  stepKind `db:"kind" json:"kind,omitempty"`
	ID string   `db:"id" json:"id,omitempty"`
}

type rootData struct {
	ID   string   `db:"id"`
	K    stepKind `db:"kind"`
	VID  string   `db:"via_id"`
	Term termData `db:"payload"`
}

type stepKind int

const (
	nonstep = stepKind(iota)
	proc
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
	X    string   `json:"x"`
	Cont termData `json:"cont"`
}

type sendData struct {
	A string `json:"a"`
	B string `json:"b"`
}

type recvData struct {
	X    string   `json:"x"`
	Y    string   `json:"y"`
	Cont termData `json:"cont"`
}

type labData struct {
	C     string `json:"c"`
	Label string `json:"label"`
}

type caseData struct {
	Z     string              `json:"z"`
	Conts map[string]termData `json:"conts"`
}

type termKind int

const (
	nonterm = termKind(iota)
	fwd
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
	DataToTerms    func([]termData) ([]Term, error)
	DataFromTerms  func([]Term) []termData
	DataToValues   func([]termData) ([]Value, error)
	DataFromValues func([]Value) []termData
	DataToConts    func([]termData) ([]Continuation, error)
	DataFromConts  func([]Continuation) []termData
)

func dataFromRoot(r root) (*rootData, error) {
	if r == nil {
		return nil, nil
	}
	switch root := r.(type) {
	case MsgRoot:
		return &rootData{
			K:    msg,
			ID:   root.ID.String(),
			VID:  root.VID.String(),
			Term: dataFromValue(root.Val),
		}, nil
	case SrvRoot:
		return &rootData{
			K:    srv,
			ID:   root.ID.String(),
			VID:  root.VID.String(),
			Term: dataFromCont(root.Cont),
		}, nil
	default:
		panic(ErrUnexpectedStep(root))
	}
}

func dataToRoot(dto *rootData) (root, error) {
	if dto == nil {
		return nil, nil
	}
	rootID, err := id.StringToID(dto.ID)
	if err != nil {
		return nil, err
	}
	viaID, err := id.StringToID(dto.VID)
	if err != nil {
		return nil, err
	}
	switch dto.K {
	case msg:
		val, err := dataToValue(dto.Term)
		if err != nil {
			return nil, err
		}
		return MsgRoot{ID: rootID, VID: viaID, Val: val}, nil
	case srv:
		cont, err := dataToCont(dto.Term)
		if err != nil {
			return nil, err
		}
		return SrvRoot{ID: rootID, VID: viaID, Cont: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func dataFromTerm(t Term) termData {
	switch term := t.(type) {
	case CloseSpec:
		return dataFromValue(term)
	case WaitSpec:
		return dataFromCont(term)
	case SendSpec:
		return dataFromValue(term)
	case RecvSpec:
		return dataFromCont(term)
	case LabSpec:
		return dataFromValue(term)
	case CaseSpec:
		return dataFromCont(term)
	default:
		panic(ErrUnexpectedTerm(term))
	}
}

func dataToTerm(dto termData) (Term, error) {
	switch dto.K {
	case close:
		return dataToValue(dto)
	case wait:
		return dataToCont(dto)
	case send:
		return dataToValue(dto)
	case recv:
		return dataToCont(dto)
	case lab:
		return dataToValue(dto)
	case caseK:
		return dataToCont(dto)
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromValue(v Value) termData {
	switch val := v.(type) {
	case CloseSpec:
		return termData{
			K:     close,
			Close: &closeData{val.A.String()},
		}
	case SendSpec:
		return termData{
			K:    send,
			Send: &sendData{val.A.String(), val.B.String()},
		}
	case LabSpec:
		return termData{
			K:   lab,
			Lab: &labData{val.C.String(), string(val.L)},
		}
	default:
		panic(ErrUnexpectedValue(val))
	}
}

func dataToValue(dto termData) (Value, error) {
	switch dto.K {
	case close:
		a, err := id.StringToID(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case send:
		a, err := id.StringToID(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := id.StringToID(dto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case lab:
		c, err := id.StringToID(dto.Lab.C)
		if err != nil {
			return nil, err
		}
		return LabSpec{C: c, L: state.Label(dto.Lab.Label)}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromCont(c Continuation) termData {
	switch cont := c.(type) {
	case WaitSpec:
		return termData{
			K: wait,
			Wait: &waitData{
				X:    cont.X.String(),
				Cont: dataFromTerm(cont.Cont),
			},
		}
	case RecvSpec:
		return termData{
			K: recv,
			Recv: &recvData{
				X:    cont.X.String(),
				Y:    cont.Y.String(),
				Cont: dataFromTerm(cont.Cont),
			},
		}
	case CaseSpec:
		conts := make(map[string]termData, len(cont.Conts))
		for l, t := range cont.Conts {
			conts[string(l)] = dataFromTerm(t)
		}
		return termData{
			K: caseK,
			Case: &caseData{
				Z:     cont.Z.String(),
				Conts: conts,
			},
		}
	default:
		panic(ErrUnexpectedCont(cont))
	}
}

func dataToCont(dto termData) (Continuation, error) {
	switch dto.K {
	case wait:
		x, err := id.StringToID(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case recv:
		x, err := id.StringToID(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := id.StringToID(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	case caseK:
		z, err := id.StringToID(dto.Case.Z)
		if err != nil {
			return nil, err
		}
		branches := make(map[state.Label]Term, len(dto.Case.Conts))
		for l, t := range dto.Case.Conts {
			branch, err := dataToTerm(t)
			if err != nil {
				return nil, err
			}
			branches[state.Label(l)] = branch
		}
		return CaseSpec{Z: z, Conts: branches}, nil
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
