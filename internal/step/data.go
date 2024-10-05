package step

import (
	"database/sql"
	"fmt"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"
	"smecalculus/rolevod/lib/sym"

	"smecalculus/rolevod/internal/chnl"
	"smecalculus/rolevod/internal/state"
)

type refData struct {
	K  stepKind `db:"kind" json:"kind,omitempty"`
	ID string   `db:"id" json:"id,omitempty"`
}

type rootData struct {
	ID   string         `db:"id"`
	K    stepKind       `db:"kind"`
	PID  sql.NullString `db:"pid"`
	VID  sql.NullString `db:"vid"`
	Term termData       `db:"term"`
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
	Caze  *cazeData  `json:"case,omitempty"`
	CTA   *ctaData   `json:"cta,omitempty"`
}

type closeData struct {
	A core.PlaceholderDTO `json:"a"`
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

type cazeData struct {
	Z     string              `json:"z"`
	Conts map[string]termData `json:"conts"`
}

type ctaData struct {
	Seat string `json:"seat"`
	Key  string `json:"key"`
}

type termKind int

const (
	nonterm = termKind(iota)
	close
	wait
	send
	recv
	lab
	caze
	cta
	link
	spawn
	fwd
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

func dataFromRoot(r Root) (*rootData, error) {
	if r == nil {
		return nil, nil
	}
	switch root := r.(type) {
	case ProcRoot:
		pid := sql.NullString{String: root.PID.String(), Valid: !root.PID.IsEmpty()}
		return &rootData{
			K:    proc,
			ID:   root.ID.String(),
			PID:  pid,
			Term: dataFromTerm(root.Term),
		}, nil
	case MsgRoot:
		vid := sql.NullString{String: root.VID.String(), Valid: !root.VID.IsEmpty()}
		return &rootData{
			K:    msg,
			ID:   root.ID.String(),
			VID:  vid,
			Term: dataFromValue(root.Val),
		}, nil
	case SrvRoot:
		pid := sql.NullString{String: root.PID.String(), Valid: !root.PID.IsEmpty()}
		vid := sql.NullString{String: root.VID.String(), Valid: !root.VID.IsEmpty()}
		return &rootData{
			K:    srv,
			ID:   root.ID.String(),
			PID:  pid,
			VID:  vid,
			Term: dataFromCont(root.Cont),
		}, nil
	default:
		panic(ErrUnexpectedStep(root))
	}
}

func dataToRoot(dto *rootData) (Root, error) {
	if dto == nil {
		return nil, nil
	}
	rid, err := id.StringToID(dto.ID)
	if err != nil {
		return nil, err
	}
	var pid chnl.ID
	if dto.PID.Valid {
		pid, err = id.StringToID(dto.PID.String)
		if err != nil {
			return nil, err
		}
	}
	var vid chnl.ID
	if dto.VID.Valid {
		vid, err = id.StringToID(dto.VID.String)
		if err != nil {
			return nil, err
		}
	}
	switch dto.K {
	case proc:
		term, err := dataToTerm(dto.Term)
		if err != nil {
			return nil, err
		}
		return ProcRoot{ID: rid, PID: pid, Term: term}, nil
	case msg:
		val, err := dataToValue(dto.Term)
		if err != nil {
			return nil, err
		}
		return MsgRoot{ID: rid, VID: vid, Val: val}, nil
	case srv:
		cont, err := dataToCont(dto.Term)
		if err != nil {
			return nil, err
		}
		return SrvRoot{ID: rid, PID: pid, VID: vid, Cont: cont}, nil
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
	case CTASpec:
		return termData{
			K: cta,
			CTA: &ctaData{
				Seat: term.Seat.Name(),
				Key:  term.Key.String(),
			},
		}
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
	case caze:
		return dataToCont(dto)
	case cta:
		key, err := ak.StringToAK(dto.CTA.Key)
		if err != nil {
			return nil, err
		}
		return CTASpec{
			Seat: sym.StringToSym(dto.CTA.Seat),
			Key:  key,
		}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromValue(v Value) termData {
	switch val := v.(type) {
	case CloseSpec:
		return termData{
			K:     close,
			Close: &closeData{core.MsgFromPH(val.A)},
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
		a, err := core.MsgToPH(dto.Close.A)
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
			K: caze,
			Caze: &cazeData{
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
	case caze:
		z, err := id.StringToID(dto.Caze.Z)
		if err != nil {
			return nil, err
		}
		branches := make(map[state.Label]Term, len(dto.Caze.Conts))
		for l, t := range dto.Caze.Conts {
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
