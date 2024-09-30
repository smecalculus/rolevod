package sym

import (
	"strings"
)

type ADT string

func New(name string) ADT {
	return ADT(name)
}

func (ns ADT) New(name string) ADT {
	return ADT(strings.Join([]string{string(ns), name}, "."))
}

func Ident(s ADT) ADT {
	return s
}

func StringToSym(s string) ADT {
	return ADT(s)
}

func StringFromSym(s ADT) string {
	return string(s)
}
