package ph

import (
	"fmt"
)

type ADT interface {
	PH()
}

func errUnexpectedType(ph ADT) error {
	return fmt.Errorf("unexpected placeholder type: %T", ph)
}
