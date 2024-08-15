package core

import (
	"github.com/rs/xid"
)

type Entity interface{}

type ID[T Entity] xid.ID

func New[T Entity]() ID[T] {
	return ID[T](xid.New())
}

func ToString[T Entity](id ID[T]) string {
	return xid.ID(id).String()
}

func FromString[T Entity](sid string) (ID[T], error) {
	cid, err := xid.FromString(sid)
	if err != nil {
		return ID[T]{}, err
	}
	return ID[T](cid), nil
}
