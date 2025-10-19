package svg

import "github.com/abzicht/svgocode/svgocode/math64"

type Translate math64.VectorF2

func (t *Translate) apply(p math64.VectorF2) math64.VectorF2 {
	return p.Add(t)
}
func (t *Translate) Apply(p ...math64.VectorF2) []math64.VectorF2 {
	for i, _ := range p {
		p[i] = t.apply(p[i])
	}
	return p
}
