package math64

import (
	"github.com/abzicht/svgocode/llog"
)

type MatrixF []Float
type MatrixF2 MatrixF
type MatrixF3 MatrixF
type MatrixF4 MatrixF

// Stored row-major, i.e.,
// [ m00, m01, m02, m03,
//   m10, m11, m12, m13,
//   m20, m21, m22, m23,
//   m30, m31, m32, m33 ]

func NewMatrixF2(values [4]Float) MatrixF2 {
	var m MatrixF2 = values[:]
	return m
}
func NewMatrixF3(values [9]Float) MatrixF3 {
	var m MatrixF3 = values[:]
	return m
}
func NewMatrixF4(values [16]Float) MatrixF4 {
	var m MatrixF4 = values[:]
	return m
}

func equal(m1, m2 MatrixF) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i, _ := range m1 {
		if m1[i] != m2[i] {
			return false
		}
	}
	return true
}

func (m MatrixF2) Equal(m2 MatrixF2) bool {
	return equal(MatrixF(m), MatrixF(m2))
}

func (m MatrixF3) Equal(m2 MatrixF3) bool {
	return equal(MatrixF(m), MatrixF(m2))
}

func (m MatrixF4) Equal(m2 MatrixF4) bool {
	return equal(MatrixF(m), MatrixF(m2))
}

func MatrixF2Identity() MatrixF2 {
	return MatrixF2{
		1, 0,
		0, 1,
	}
}

func MatrixF3Identity() MatrixF3 {
	return MatrixF3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

func MatrixF4Identity() MatrixF4 {
	return MatrixF4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Matrix product, for square matrices
func mProduct(m1, m2 MatrixF, dim int) MatrixF3 {
	if len(m1) != dim*dim || len(m2) != len(m1) {
		llog.Panicf("Failed to create product of matrices. They do not match the given dimension (%d)", dim)
	}
	m3 := make([]Float, len(m1))
	for row := 0; row < dim; row++ {
		for col := 0; col < dim; col++ {
			var sum Float = 0.0
			for k := 0; k < dim; k++ {
				sum += m1[row*dim+k] * m2[k*dim+col]
			}
			m3[row*dim+col] = sum
		}
	}
	return m3
}

func (m MatrixF2) MProduct(m2 MatrixF2) MatrixF2 {
	return MatrixF2(mProduct(MatrixF(m), MatrixF(m2), 2))
}

func (m MatrixF3) MProduct(m2 MatrixF3) MatrixF3 {
	return MatrixF3(mProduct(MatrixF(m), MatrixF(m2), 3))
}

func (m MatrixF4) MProduct(m2 MatrixF4) MatrixF4 {
	return MatrixF4(mProduct(MatrixF(m), MatrixF(m2), 4))
}

func mvProduct(matrix MatrixF, vector []Float) []Float {
	n := int(len(vector))
	if len(matrix) != n*n {
		llog.Panicf("Matrix must be %dx%d for a vector of length %d", n, n, n)
	}

	result := make([]Float, n)
	for row := 0; row < n; row++ {
		var sum Float = 0.0
		for col := 0; col < n; col++ {
			sum += matrix[row*n+col] * vector[col]
		}
		result[row] = sum
	}
	return result
}

func (m MatrixF2) VProduct(v VectorF2) VectorF2 {
	result := mvProduct(MatrixF(m), []Float{v.X, v.Y})
	return VectorF2{X: result[0], Y: result[1]}
}

func (m MatrixF3) VProduct(v VectorF3) VectorF3 {
	result := mvProduct(MatrixF(m), []Float{v.X, v.Y, v.Z})
	return VectorF3{X: result[0], Y: result[1], Z: result[2]}
}

func (m MatrixF4) VProduct(v VectorF4) VectorF4 {
	result := mvProduct(MatrixF(m), []Float{v.X, v.Y, v.Z, v.W})
	return VectorF4{X: result[0], Y: result[1], Z: result[2], W: result[3]}
}
