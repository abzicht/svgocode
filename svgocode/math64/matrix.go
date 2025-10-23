package math64

import (
	"fmt"

	"github.com/abzicht/svgocode/llog"
)

type MatrixF3 struct {
	A VectorF3 // A.X B.X C.X
	B VectorF3 // A.Y B.Y C.Y
	C VectorF3 // A.Z B.Z C.Z
}

type MatrixF2 struct {
	A VectorF2
	B VectorF2
}

type MatrixF32 struct {
	A VectorF3
	B VectorF3
}

// Creates a new matrix from the given values
// Missing values are initialized with 0.
// More than 9 values result in a panic.
// Fills values in the order
// 1 2 3
// 4 5 6
// 7 8 9.
func NewMatrixF3(values ...Float) MatrixF3 {
	if len(values) > 9 {
		llog.Panicf("Cannot create 3x3 matrix from %d values (too many)", len(values))
	}
	m := MatrixF3{}
	for i, val := range values {
		switch i {
		case 0:
			m.A.X = val
		case 1:
			m.B.X = val
		case 2:
			m.C.X = val
		case 3:
			m.A.Y = val
		case 4:
			m.B.Y = val
		case 5:
			m.C.Y = val
		case 6:
			m.A.Z = val
		case 7:
			m.B.Z = val
		case 8:
			m.C.Z = val
		}
	}
	return m
}

func (m MatrixF2) Equal(m2 MatrixF2) bool {
	return m.A.Equal(m2.A) && m.B.Equal(m2.B)
}

func (m MatrixF3) Equal(m2 MatrixF3) bool {
	return m.A.Equal(m2.A) && m.B.Equal(m2.B) && m.C.Equal(m2.C)
}

func (m MatrixF32) Equal(m2 MatrixF32) bool {
	return m.A.Equal(m2.A) && m.B.Equal(m2.B)
}

func (m MatrixF3) Identity() MatrixF3 {
	return MatrixF3{
		A: VectorF3{X: 1, Y: 0, Z: 0},
		B: VectorF3{X: 0, Y: 1, Z: 0},
		C: VectorF3{X: 0, Y: 0, Z: 1},
	}
}

func (m MatrixF3) MProduct(m2 MatrixF3) MatrixF3 {
	ax := m.A.X*m2.A.X + m.B.X*m2.A.Y + m.C.X*m2.A.Z
	ay := m.A.Y*m2.A.X + m.B.Y*m2.A.Y + m.C.Y*m2.A.Z
	az := m.A.Z*m2.A.X + m.B.Z*m2.A.Y + m.C.Z*m2.A.Z

	bx := m.A.X*m2.B.X + m.B.X*m2.B.Y + m.C.X*m2.B.Z
	by := m.A.Y*m2.B.X + m.B.Y*m2.B.Y + m.C.Y*m2.B.Z
	bz := m.A.Z*m2.B.X + m.B.Z*m2.B.Y + m.C.Z*m2.B.Z

	cx := m.A.X*m2.C.X + m.B.X*m2.C.Y + m.C.X*m2.C.Z
	cy := m.A.Y*m2.C.X + m.B.Y*m2.C.Y + m.C.Y*m2.C.Z
	cz := m.A.Z*m2.C.X + m.B.Z*m2.C.Y + m.C.Z*m2.C.Z
	return MatrixF3{
		A: VectorF3{X: ax, Y: ay, Z: az},
		B: VectorF3{X: bx, Y: by, Z: bz},
		C: VectorF3{X: cx, Y: cy, Z: cz},
	}
}

func (m MatrixF3) VProductF3(v VectorF3) VectorF3 {
	x := v.X*m.A.X + v.Y*m.B.X + v.Z*m.C.X
	y := v.X*m.A.Y + v.Y*m.B.Y + v.Z*m.C.Y
	z := v.X*m.A.Z + v.Y*m.B.Z + v.Z*m.C.Z
	return VectorF3{X: x, Y: y, Z: z}
}

func (m MatrixF3) VProductF2(v VectorF2) VectorF2 {
	v3 := m.VProductF3(VectorF3{X: v.X, Y: v.Y, Z: 1})
	return VectorF2{X: v3.X, Y: v3.Y}
}

func (m MatrixF2) String() string {
	return fmt.Sprintf("(%f %f, %f %f)", m.A.X, m.A.Y, m.B.X, m.B.Y)
}

func (m MatrixF3) String() string {
	return fmt.Sprintf("(%f %f %f, %f %f %f, %f %f %f)", m.A.X, m.A.Y, m.A.Z, m.B.X, m.B.Y, m.B.Z, m.C.X, m.C.Y, m.C.Z)
}

func (m MatrixF32) String() string {
	return fmt.Sprintf("(%f %f %f, %f %f %f)", m.A.X, m.A.Y, m.A.Z, m.B.X, m.B.Y, m.B.Z)
}
