package id

import (
	"github.com/rs/xid"
)

type ADT xid.ID

func New() ADT {
	return ADT(xid.New())
}

func Ident(id ADT) ADT {
	return id
}

func (id ADT) String() string {
	return xid.ID(id).String()
}

func StringTo(id ADT) string {
	return xid.ID(id).String()
}

func StringFrom(s string) (ADT, error) {
	xid, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(xid), nil
}
