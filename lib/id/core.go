package id

import (
	"errors"

	"github.com/rs/xid"
)

type Identifiable interface {
	Ident() ADT
}

type ADT xid.ID

func (ADT) PH() {}

func New() ADT {
	return ADT(xid.New())
}

func Empty() ADT {
	return ADT(xid.NilID())
}

func (id ADT) IsEmpty() bool {
	return xid.ID(id).IsZero()
}

func Ident(id ADT) ADT {
	return id
}

func StringToID(s string) (ADT, error) {
	xid, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(xid), nil
}

func StringFromID(id ADT) string {
	return xid.ID(id).String()
}

func StringFromID2(id *ADT) *string {
	if id == nil {
		return nil
	}
	s := xid.ID(*id).String()
	return &s
}

func (id ADT) String() string {
	return xid.ID(id).String()
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend StringToID
// goverter:extend StringFromID
var (
	StringsToIDs   func([]string) ([]ADT, error)
	StringsFromIDs func([]ADT) []string
)

var (
	ErrEmpty = errors.New("empty id")
)
