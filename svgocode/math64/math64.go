package math64

// mm
type Float float64

// mm/s
type Speed Float

type VectorF3 struct {
	X Float
	Y Float
	Z Float
}

type VectorF2 struct {
	X Float
	Y Float
}

func (v VectorF2) Equal(v2 VectorF2) bool {
	return v.X == v2.X && v.Y == v2.Y
}

func (v VectorF3) Equal(v2 VectorF3) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}
