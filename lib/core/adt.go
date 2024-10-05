package core

import (
	"fmt"
)

type Placeholder interface {
	PH()
}

func errUnexpectedPHType(ph Placeholder) error {
	return fmt.Errorf("unexpected placeholder type: %T", ph)
}
