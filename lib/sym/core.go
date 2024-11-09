package sym

import (
	"strings"
)

type Symbolizable interface {
	Sym() ADT
}

type ADT string

func (ADT) PH() {}

func New(name string) ADT {
	return ADT(name)
}

func (ns ADT) New(name string) ADT {
	return ADT(strings.Join([]string{string(ns), name}, sep))
}

func (s ADT) Name() string {
	sym := string(s)
	return sym[strings.LastIndex(sym, sep)+1:]
}

func (s ADT) NS() ADT {
	sym := string(s)
	return ADT(sym[0:strings.LastIndex(sym, sep)])
}

func ConvertToSame(a ADT) ADT {
	return a
}

func CovertFromString(s string) ADT {
	return ADT(s)
}

func ConvertToString(a ADT) string {
	return string(a)
}

const (
	sep = "."
)
