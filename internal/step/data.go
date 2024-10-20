package step

import (
	"database/sql"
	"fmt"

	"smecalculus/rolevod/lib/ak"
	"smecalculus/rolevod/lib/core"
	"smecalculus/rolevod/lib/id"

	"smecalculus/rolevod/internal/chnl"
)

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
	K     termKind   `json:"k"`
	Close *closeData `json:"close,omitempty"`
	Wait  *waitData  `json:"wait,omitempty"`
	Send  *sendData  `json:"send,omitempty"`
	Recv  *recvData  `json:"recv,omitempty"`
	Lab   *labData   `json:"lab,omitempty"`
	Case  *caseData  `json:"case,omitempty"`
	Fwd   *fwdData   `json:"fwd,omitempty"`
	CTA   *ctaData   `json:"cta,omitempty"`
}

type closeData struct {
	A core.PlaceholderDTO `json:"a"`
}

type waitData struct {
	X    core.PlaceholderDTO `json:"x"`
	Cont termData            `json:"cont"`
}

type sendData struct {
	A core.PlaceholderDTO `json:"a"`
	B core.PlaceholderDTO `json:"b"`
}

type recvData struct {
	X    core.PlaceholderDTO `json:"x"`
	Y    core.PlaceholderDTO `json:"y"`
	Cont termData            `json:"cont"`
}

type labData struct {
	C     core.PlaceholderDTO `json:"c"`
	Label string              `json:"l"`
}

type caseData struct {
	Z     core.PlaceholderDTO `json:"z"`
	Conts map[string]termData `json:"conts"`
}

type fwdData struct {
	C core.PlaceholderDTO `json:"c"`
	D core.PlaceholderDTO `json:"d"`
}

type ctaData struct {
	AK   string `json:"ak"`
	Seat string `json:"seat_id"`
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
	DataFromTerms  func([]Term) ([]termData, error)
	DataToValues   func([]termData) ([]Value, error)
	DataFromValues func([]Value) []termData
	DataToConts    func([]termData) ([]Continuation, error)
	DataFromConts  func([]Continuation) ([]termData, error)
)

func dataFromRoot(r Root) (*rootData, error) {
	if r == nil {
		return nil, nil
	}
	switch root := r.(type) {
	case ProcRoot:
		pid := sql.NullString{String: root.PID.String(), Valid: !root.PID.IsEmpty()}
		term, err := dataFromTerm(root.Term)
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:    proc,
			ID:   root.ID.String(),
			PID:  pid,
			Term: term,
		}, nil
	case MsgRoot:
		pid := sql.NullString{String: root.PID.String(), Valid: !root.PID.IsEmpty()}
		vid := sql.NullString{String: root.VID.String(), Valid: !root.VID.IsEmpty()}
		return &rootData{
			K:    msg,
			ID:   root.ID.String(),
			PID:  pid,
			VID:  vid,
			Term: dataFromValue(root.Val),
		}, nil
	case SrvRoot:
		pid := sql.NullString{String: root.PID.String(), Valid: !root.PID.IsEmpty()}
		vid := sql.NullString{String: root.VID.String(), Valid: !root.VID.IsEmpty()}
		term, err := dataFromCont(root.Cont)
		if err != nil {
			return nil, err
		}
		return &rootData{
			K:    srv,
			ID:   root.ID.String(),
			PID:  pid,
			VID:  vid,
			Term: term,
		}, nil
	default:
		panic(ErrRootTypeUnexpected(root))
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
		return MsgRoot{ID: rid, PID: pid, VID: vid, Val: val}, nil
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

func dataFromTerm(t Term) (termData, error) {
	switch term := t.(type) {
	case CloseSpec:
		return dataFromValue(term), nil
	case WaitSpec:
		return dataFromCont(term)
	case SendSpec:
		return dataFromValue(term), nil
	case RecvSpec:
		return dataFromCont(term)
	case LabSpec:
		return dataFromValue(term), nil
	case CaseSpec:
		return dataFromCont(term)
	case FwdSpec:
		return dataFromValue(term), nil
	case CTASpec:
		return termData{
			K: cta,
			CTA: &ctaData{
				Seat: term.Seat.String(),
				AK:   term.AK.String(),
			},
		}, nil
	default:
		panic(ErrTermTypeUnexpected(term))
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
	case fwd:
		return dataToValue(dto)
	case cta:
		key, err := ak.StringToAK(dto.CTA.AK)
		if err != nil {
			return nil, err
		}
		seat, err := id.StringToID(dto.CTA.Seat)
		if err != nil {
			return nil, err
		}
		return CTASpec{Seat: seat, AK: key}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromValue(v Value) termData {
	switch val := v.(type) {
	case CloseSpec:
		return termData{
			K:     close,
			Close: &closeData{core.DTOFromPH(val.A)},
		}
	case SendSpec:
		return termData{
			K:    send,
			Send: &sendData{core.DTOFromPH(val.A), core.DTOFromPH(val.B)},
		}
	case LabSpec:
		return termData{
			K:   lab,
			Lab: &labData{core.DTOFromPH(val.C), string(val.L)},
		}
	case FwdSpec:
		return termData{
			K: fwd,
			Fwd: &fwdData{
				C: core.DTOFromPH(val.C),
				D: core.DTOFromPH(val.D),
			},
		}
	default:
		panic(ErrValTypeUnexpected(val))
	}
}

func dataToValue(dto termData) (Value, error) {
	switch dto.K {
	case close:
		a, err := core.DTOToPH(dto.Close.A)
		if err != nil {
			return nil, err
		}
		return CloseSpec{A: a}, nil
	case send:
		a, err := core.DTOToPH(dto.Send.A)
		if err != nil {
			return nil, err
		}
		b, err := core.DTOToPH(dto.Send.B)
		if err != nil {
			return nil, err
		}
		return SendSpec{A: a, B: b}, nil
	case lab:
		c, err := core.DTOToPH(dto.Lab.C)
		if err != nil {
			return nil, err
		}
		return LabSpec{C: c, L: core.Label(dto.Lab.Label)}, nil
	case fwd:
		c, err := core.DTOToPH(dto.Fwd.C)
		if err != nil {
			return nil, err
		}
		d, err := core.DTOToPH(dto.Fwd.D)
		if err != nil {
			return nil, err
		}
		return FwdSpec{C: c, D: d}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func dataFromCont(c Continuation) (termData, error) {
	switch cont := c.(type) {
	case WaitSpec:
		dto, err := dataFromTerm(cont.Cont)
		if err != nil {
			return termData{}, err
		}
		return termData{
			K: wait,
			Wait: &waitData{
				X:    core.DTOFromPH(cont.X),
				Cont: dto,
			},
		}, nil
	case RecvSpec:
		dto, err := dataFromTerm(cont.Cont)
		if err != nil {
			return termData{}, err
		}
		return termData{
			K: recv,
			Recv: &recvData{
				X:    core.DTOFromPH(cont.X),
				Y:    core.DTOFromPH(cont.Y),
				Cont: dto,
			},
		}, nil
	case CaseSpec:
		conts := make(map[string]termData, len(cont.Conts))
		for l, t := range cont.Conts {
			dto, err := dataFromTerm(t)
			if err != nil {
				return termData{}, err
			}
			conts[string(l)] = dto
		}
		return termData{
			K: caze,
			Case: &caseData{
				Z:     core.DTOFromPH(cont.Z),
				Conts: conts,
			},
		}, nil
	case FwdSpec:
		return termData{
			K: fwd,
			Fwd: &fwdData{
				C: core.DTOFromPH(cont.C),
				D: core.DTOFromPH(cont.D),
			},
		}, nil
	default:
		panic(ErrContTypeUnexpected(cont))
	}
}

func dataToCont(dto termData) (Continuation, error) {
	switch dto.K {
	case wait:
		x, err := core.DTOToPH(dto.Wait.X)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Wait.Cont)
		if err != nil {
			return nil, err
		}
		return WaitSpec{X: x, Cont: cont}, nil
	case recv:
		x, err := core.DTOToPH(dto.Recv.X)
		if err != nil {
			return nil, err
		}
		y, err := core.DTOToPH(dto.Recv.Y)
		if err != nil {
			return nil, err
		}
		cont, err := dataToTerm(dto.Recv.Cont)
		if err != nil {
			return nil, err
		}
		return RecvSpec{X: x, Y: y, Cont: cont}, nil
	case caze:
		z, err := core.DTOToPH(dto.Case.Z)
		if err != nil {
			return nil, err
		}
		branches := make(map[core.Label]Term, len(dto.Case.Conts))
		for l, t := range dto.Case.Conts {
			branch, err := dataToTerm(t)
			if err != nil {
				return nil, err
			}
			branches[core.Label(l)] = branch
		}
		return CaseSpec{Z: z, Conts: branches}, nil
	case fwd:
		c, err := core.DTOToPH(dto.Fwd.C)
		if err != nil {
			return nil, err
		}
		d, err := core.DTOToPH(dto.Fwd.D)
		if err != nil {
			return nil, err
		}
		return FwdSpec{C: c, D: d}, nil
	default:
		panic(errUnexpectedTermKind(dto.K))
	}
}

func errUnexpectedTermKind(k termKind) error {
	return fmt.Errorf("unexpected term kind: %v", k)
}

func errUnexpectedStepKind(k stepKind) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
