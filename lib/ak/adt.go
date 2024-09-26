package ak

import (
	"github.com/rs/xid"
)

type ADT xid.ID

func New() ADT {
	return ADT(xid.New())
}

func Ident(ak ADT) ADT {
	return ak
}

func (ak ADT) String() string {
	return xid.ID(ak).String()
}

func StringFrom(ak ADT) string {
	return xid.ID(ak).String()
}

func StringTo(s string) (ADT, error) {
	xid, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(xid), nil
}
