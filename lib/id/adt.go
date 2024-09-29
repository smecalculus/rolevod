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
