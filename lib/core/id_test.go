package core

import (
	"fmt"
	"testing"
)

type Foo Entity

func TestHelloName(t *testing.T) {
	fmt.Print(ToString(New[Foo]()))
}
