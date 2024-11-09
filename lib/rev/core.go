package rev

import (
	"time"
)

type ADT int64

func New() ADT {
	return ADT(time.Now().Unix())
}

func ConvertToInt(a ADT) int64 {
	return int64(a)
}

func ConvertFromInt(i int64) ADT {
	return ADT(i)
}
