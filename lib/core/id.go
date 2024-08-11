package core

import (
	"github.com/rs/xid"
)

type Kind interface{}

type Id[T Kind] xid.ID

func New[T Kind]() Id[T] {
	return Id[T](xid.New())
}
