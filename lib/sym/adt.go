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

func Ident(s ADT) ADT {
	return s
}

func StringToSym(s string) ADT {
	return ADT(s)
}

func StringFromSym(s ADT) string {
	return string(s)
}

const (
	sep = "."
)
