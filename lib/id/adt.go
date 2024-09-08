package id

import (
	"github.com/rs/xid"
)

type ADT[T any] xid.ID

func New[T any]() ADT[T] {
	return ADT[T](xid.New())
}

func (id ADT[T]) String() string {
	return xid.ID(id).String()
}

func String[T any](sid string) (ADT[T], error) {
	xid, err := xid.FromString(sid)
	if err != nil {
		return ADT[T]{}, err
	}
	return ADT[T](xid), nil
}
