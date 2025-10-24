package math64

import "fmt"

type VectorT2[T comparable] struct {
	X T
	Y T
}

type VectorT3[T comparable] struct {
	X T
	Y T
	Z T
}

type VectorT4[T comparable] struct {
	X T
	Y T
	Z T
	W T
}

type VectorF2 VectorT2[Float]
type VectorF3 VectorT3[Float]
type VectorF4 VectorT4[Float]

func (v VectorF2) Equal(v2 VectorF2) bool {
	return v.X == v2.X && v.Y == v2.Y
}

func (v VectorF3) Equal(v2 VectorF3) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

func (v VectorF4) Equal(v2 VectorF4) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z && v.W == v2.W
}

func (v VectorF2) Add(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v VectorF3) Add(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

func (v VectorF4) Add(v2 VectorF4) VectorF4 {
	return VectorF4{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z, W: v.W + v2.W}
}

func (v VectorF2) Min(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X.Min(v2.X), Y: v.Y.Min(v2.Y)}
}

func (v VectorF3) Min(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X.Min(v2.X), Y: v.Y.Min(v2.Y), Z: v.Z.Min(v2.Z)}
}

func (v VectorF4) Min(v2 VectorF4) VectorF4 {
	return VectorF4{X: v.X.Min(v2.X), Y: v.Y.Min(v2.Y), Z: v.Z.Min(v2.Z), W: v.W.Min(v2.W)}
}

func (v VectorF2) Max(v2 VectorF2) VectorF2 {
	return VectorF2{X: v.X.Max(v2.X), Y: v.Y.Max(v2.Y)}
}

func (v VectorF3) Max(v2 VectorF3) VectorF3 {
	return VectorF3{X: v.X.Max(v2.X), Y: v.Y.Max(v2.Y), Z: v.Z.Max(v2.Z)}
}

func (v VectorF4) Max(v2 VectorF4) VectorF4 {
	return VectorF4{X: v.X.Max(v2.X), Y: v.Y.Max(v2.Y), Z: v.Z.Max(v2.Z), W: v.W.Max(v2.W)}
}

func (v VectorF2) String() string {
	return fmt.Sprintf("%fx %fy", v.X, v.Y)
}

func (v VectorF3) String() string {
	return fmt.Sprintf("%fx %fy %fz", v.X, v.Y, v.Z)
}

func (v VectorF4) String() string {
	return fmt.Sprintf("%fx %fy %fz %fw", v.X, v.Y, v.Z, v.W)
}
