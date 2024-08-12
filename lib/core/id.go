package core

import (
	"github.com/rs/xid"
)

type Kind interface{}

type Id[T Kind] xid.ID

func New[T Kind]() Id[T] {
	return Id[T](xid.New())
}

func ToString[T Kind](id Id[T]) string {
	return xid.ID(id).String()
}

func FromString[T Kind](sid string) (Id[T], error) {
	cid, err := xid.FromString(sid)
	if err != nil {
		return Id[T]{}, err
	}
	return Id[T](cid), nil
}
