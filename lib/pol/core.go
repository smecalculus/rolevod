package pol

type ADT int8

const (
	Pos  = ADT(+1)
	Zero = ADT(0)
	Neg  = ADT(-1)
)

type Polarizable interface {
	Pol() ADT
}
