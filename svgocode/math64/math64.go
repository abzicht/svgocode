package math64

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

// mm/s
type Speed Float
