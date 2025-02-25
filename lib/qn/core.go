package qn

import (
	"strings"
)

type ADT string

type Qualifiable interface {
	Name() ADT
}

func New(qn string) ADT {
	return ADT(qn)
}

func (qn ADT) New(sn string) ADT {
	return ADT(strings.Join([]string{string(qn), sn}, sep))
}

// short name
func (qn ADT) SN() string {
	str := string(qn)
	return str[strings.LastIndex(str, sep)+1:]
}

// namespace
func (qn ADT) NS() ADT {
	str := string(qn)
	return ADT(str[0:strings.LastIndex(str, sep)])
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
