package math64

import "fmt"

// mm
type Float float64

func (f Float) Min(f2 Float) Float {
	if f < f2 {
		return f
	}
	return f2
}

func (f Float) Max(f2 Float) Float {
	if f > f2 {
		return f
	}
	return f2
}

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

func (v VectorF2) Add(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v VectorF3) Add(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

func (v VectorF2) Min(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X.Min(v2.X), Y: v.Y.Min(v2.Y)}
}

func (v VectorF3) Min(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X.Min(v2.X), Y: v.Y.Min(v2.Y), Z: v.Z.Min(v2.Z)}
}

func (v VectorF2) Max(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X.Max(v2.X), Y: v.Y.Max(v2.Y)}
}

func (v VectorF3) Max(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X.Max(v2.X), Y: v.Y.Max(v2.Y), Z: v.Z.Max(v2.Z)}
}

func (v VectorF2) String() string {
	return fmt.Sprintf("%fx %fy", v.X, v.Y)
}

func (v VectorF3) String() string {
	return fmt.Sprintf("%fx %fy %fz", v.X, v.Y, v.Z)
}

// mm/s
type Speed Float
