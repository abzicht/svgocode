package math64

import "testing"

func TestMProduct(t *testing.T) {
	m1 := NewMatrixF3([9]Float{0, 1, 2, 3, 4, 5, 6, 7, 8})
	m2 := MatrixF3Identity()
	m3_expected := m1

	m3 := m1.MProduct(m2)

	if !m3.Equal(m3_expected) {
		t.Errorf("Matrix product does not match. Expected:\n%s\nGot:\n%s", m3_expected.String(), m3.String())
	}

	m2 = NewMatrixF3([9]Float{-7, 3, 2, -1, 4, 6, 8, -5, 0})
	m3_expected = NewMatrixF3([9]Float{15, -6, 6, 15, 0, 30, 15, 6, 54})

	m3 = m1.MProduct(m2)
	if !m3.Equal(m3_expected) {
		t.Errorf("Matrix product does not match. Expected:\n%s\nGot:\n%s", m3_expected.String(), m3.String())
	}
}
