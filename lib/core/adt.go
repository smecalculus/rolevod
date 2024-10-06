package core

import (
	"fmt"
)

type Label string

type Placeholder interface {
	PH()
}

func errUnexpectedPHType(ph Placeholder) error {
	return fmt.Errorf("unexpected placeholder type: %T", ph)
}
