package svgtransform

import (
	"github.com/abzicht/svgocode/svgocode/math64"
)

type Mirror struct {
	X      bool            // Mirror the X axis
	Y      bool            // Mirror the Y axis
	Center math64.VectorF2 // The point that other points are being mirrored around.
}

func NewMirror(x, y bool, center math64.VectorF2) *Mirror {
	m := new(Mirror)
	m.X = x
	m.Y = y
	m.Center = center
	return m
}

func (m *Mirror) ToMatrix() *TransformMatrix {
	// translate by offset, scale with -1, and translate back
	offset := m.Center
	scale := math64.VectorF2{X: 1, Y: 1}
	if m.X {
		scale.X = -1
	} else {
		offset.X = 0
	}
	if m.Y {
		scale.Y = -1
	} else {
		offset.Y = 0
	}
	return TransformChain{NewTranslate(offset), NewScale(scale), NewTranslate(math64.VectorF2{X: -offset.X, Y: -offset.Y})}.ToMatrix()
}

func (m *Mirror) apply(p math64.VectorF2) math64.VectorF2 {
	// p' = c + (c - p)
	if m.X {
		p.X = m.Center.X + (m.Center.X - p.X)
	}
	if m.Y {
		p.Y = m.Center.Y + (m.Center.Y - p.Y)
	}
	return p
}
func (m *Mirror) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = m.apply(p[i])
	}
	return p
}
