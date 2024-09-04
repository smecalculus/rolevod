package core

import (
	"github.com/rs/xid"
)

type ID[T any] xid.ID

func New[T any]() ID[T] {
	return ID[T](xid.New())
}

func (id ID[T]) String() string {
	return xid.ID(id).String()
}

func FromString[T any](sid string) (ID[T], error) {
	xid, err := xid.FromString(sid)
	if err != nil {
		return ID[T]{}, err
	}
	return ID[T](xid), nil
}

func toString(id xid.ID) string {
	return id.String()
}

func fromString(id string) (xid.ID, error) {
	return xid.FromString(id)
}
