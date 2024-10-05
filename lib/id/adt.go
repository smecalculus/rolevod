package id

import (
	"errors"

	"github.com/rs/xid"
)

type ADT xid.ID

func (ADT) PH() {}

func New() ADT {
	return ADT(xid.New())
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

func (id ADT) String() string {
	return xid.ID(id).String()
}

var (
	ErrEmpty = errors.New("empty id")
)
