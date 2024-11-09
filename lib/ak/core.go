package ak

import (
	"fmt"

	"github.com/rs/xid"
)

type ADT xid.ID

func New() ADT {
	return ADT(xid.New())
}

func (ak ADT) String() string {
	return xid.ID(ak).String()
}

func ConvertToSame(ak ADT) ADT {
	return ak
}

func ConvertFromString(s string) (ADT, error) {
	xid, err := xid.FromString(s)
	if err != nil {
		return ADT{}, err
	}
	return ADT(xid), nil
}

func ConvertToString(ak ADT) string {
	return xid.ID(ak).String()
}

func ErrUnexpectedKey(k ADT) error {
	return fmt.Errorf("unexpected access key: %v", k)
}
