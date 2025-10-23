package svgtransform

import "github.com/abzicht/svgocode/svgocode/math64"

type Translate struct {
	Offset math64.VectorF2
}

func NewTranslate(offset math64.VectorF2) *Translate {
	t := new(Translate)
	t.Offset = offset
	return t
}

func (t *Translate) apply(p math64.VectorF2) math64.VectorF2 {
	return p.Add(t.Offset)
}

func (t *Translate) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = t.apply(p[i])
	}
	return p
}

func (t *Translate) ToMatrix() *TransformMatrix {
	return NewTransformMatrix(math64.NewMatrixF3(1, 0, 0, 0, 1, 0, t.Offset.X, t.Offset.Y, 0))
}
